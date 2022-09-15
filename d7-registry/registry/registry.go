package registry

import (
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

// NeeRegistry is a simple register center, provide following functions.
// add a server and receive hearbeat to keep it alive.
// returns all alive servers and delete dead servers sync simultaneously
type NeeRegistry struct {
	timeout time.Duration // 保证一手性能
	mu      sync.Mutex    // protect following
	servers map[string]*ServerItem
}

type ServerItem struct {
	Addr  string
	start time.Time
}

const (
	defaultPath    = "/_neerpc_/registry"
	defaultTimeout = time.Minute * 5 // 5 分钟没鸟我(注册中心)就超时
)

// New create a registry instance with timeout setting
func New(timeout time.Duration) *NeeRegistry {
	return &NeeRegistry{
		servers: make(map[string]*ServerItem),
		timeout: timeout,
	}
}

var DefaultNeeRegister = New(defaultTimeout)

// putServer 添加服务实例，如果服务已经存在，则更新start
func (r *NeeRegistry) putServer(addr string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	s := r.servers[addr]
	if s == nil {
		r.servers[addr] = &ServerItem{
			addr,
			time.Now(),
		}
	} else {
		s.start = time.Now() // if exists, update start time to keep alive
	}
}

// aliveRegistry 返回可用的服务列表，如果存在超时的服务，则删除
func (r *NeeRegistry) aliveServer() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	var alive []string
	for addr, s := range r.servers {
		if r.timeout == 0 || s.start.Add(r.timeout).After(time.Now()) {
			alive = append(alive, addr)
		} else {
			delete(r.servers, addr)
		}
	}
	sort.Strings(alive)
	return alive
}

// ServeHTTp
// Runs at /_neerpc_/registry
func (r *NeeRegistry) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		// keep it simple, server is req.Header
		w.Header().Set("X-Neerpc-Servers", strings.Join(r.aliveServer(), ","))
	case "POST":
		// keep it simple, server is in req.Header
		addr := req.Header.Get("X-Neerpc-Server")
		if addr == "" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		r.putServer(addr)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// HandleHTTP registers an HTTP handler for NeeRegistry messages on registerPath
func (r *NeeRegistry) HandleHTTP(registryPath string) {
	http.Handle(registryPath, r)
	log.Println("rpc registry path:", registryPath)
}

func HandleHTTP() {
	DefaultNeeRegister.HandleHTTP(defaultPath)
}

// Heartbeat send a heartbeat message every once in a while
// it`s a helper function for a server to register or send heartbeat
func Heartbeat(registry, addr string, duration time.Duration) {
	if duration == 0 {
		// make sure there is enough time to send heart beat
		// before it`s removed from registry
		duration = defaultTimeout - time.Duration(1)*time.Minute
	}
	var err error
	err = sendHeartbeat(registry, addr)
	go func() {
		t := time.NewTicker(duration)
		for err == nil {
			<-t.C
			err = sendHeartbeat(registry, addr)
		}
	}()

}

func sendHeartbeat(registry, addr string) error {
	log.Println(addr, "send heart beat to registry", registry)
	httpClient := &http.Client{}
	req, _ := http.NewRequest("POST", registry, nil)
	req.Header.Set("X-Neerpc-Server", addr)
	if _, err := httpClient.Do(req); err != nil {
		return err
	}
	return nil
}

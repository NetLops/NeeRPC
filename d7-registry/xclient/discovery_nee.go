package xclient

import (
	"log"
	"net/http"
	"strings"
	"time"
)

type NeeRegistryDiscovery struct {
	*MultiServersDiscovery               // 套娃复用功能
	registry               string        // 注册中心
	timeout                time.Duration // 服务列表的过期时间
	lastUpdate             time.Time     // 带白哦从注册中心更新服务列表的时间，默认是 10s 过期，即10s 之后，需要从注册中心更新新的列表
}

const defaultUpdateTimeout = time.Second * 10

func NewNeeRegistryDiscovery(registerAddr string, timeout time.Duration) *NeeRegistryDiscovery {
	if timeout == 0 {
		timeout = defaultUpdateTimeout
	}
	d := &NeeRegistryDiscovery{
		MultiServersDiscovery: NewMultiServersDiscovery(make([]string, 0)),
		registry:              registerAddr,
		timeout:               timeout,
	}
	return d
}

func (d *NeeRegistryDiscovery) Update(servers []string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.servers = servers
	d.lastUpdate = time.Now()
	return nil
}

func (d *NeeRegistryDiscovery) Refresh() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.lastUpdate.Add(d.timeout).After(time.Now()) {
		return nil
	}

	log.Println("rpc registry: refresh servers fro registry", d.registry)
	resp, err := http.Get(d.registry)
	if err != nil {
		log.Println("rpc registry refresh err:", err)
		return err
	}
	servers := strings.Split(resp.Header.Get("X-Neerpc-Servers"), ",")
	d.servers = make([]string, 0, len(servers))
	for _, server := range servers {
		if strings.TrimSpace(server) != "" {
			d.servers = append(d.servers, strings.TrimSpace(server))
		}
	}
	d.lastUpdate = time.Now()
	return nil
}

/*
	Get 和 GetAll 和 MultiServersDiscovery 相似，唯一不同在于，NeeRegistryDiscovery 需要先调用 Refresh 确保服务没有过期
*/

func (d *NeeRegistryDiscovery) Get(mode SelectMode) (string, error) {
	if err := d.Refresh(); err != nil {
		return "", err
	}
	return d.MultiServersDiscovery.Get(mode)
}

func (d *NeeRegistryDiscovery) GetAll() ([]string, error) {
	if err := d.Refresh(); err != nil {
		return nil, err
	}
	return d.MultiServersDiscovery.GetAll()
}

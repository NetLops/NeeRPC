package d1_codec

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"neerpc/codec"
	"net"
	"reflect"
	"sync"
)

/*
1. 报文的最开始会规划固定的字节, 来协商相关的信息。
2. 第1个字节用来表示序列化方式，第2个字节表示压缩方式，
   第3-6字节表示 header 的长度，7-10 字节表示 body 的长度
*/

const MagicNumber = 0x3bef5c

type Option struct {
	MagicNumber int        // MagicNumber marks this`s neerpc request
	CodecType   codec.Type // client may choose different Codec to encode body
}

var DdefaultOption = &Option{
	MagicNumber: MagicNumber,
	CodecType:   codec.GobType,
}

// Server represents an RPC Server
type Server struct {
}

// NewServer returns a new Server.
func NewServer() *Server {
	return &Server{}
}

// Default is the default instance of *Server
var DefaultServer = NewServer()

// Accept accepts connections on the listener and servers requests
// for each incoming connection.
func (server *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc server: accept error:", err)
			return
		}
		go server.ServeConn(conn)
	}
}

// Accept accepts connections on the listener and serves requests
// for each incoming connection.
func Accept(lis net.Listener) {
	DefaultServer.Accept(lis)
}

// ServeConn runs the server on a single connection.
// ServeConn blocks, serving the connection until the client hangs up.
func (server *Server) ServeConn(conn io.ReadWriteCloser) {
	defer func() { _ = conn.Close() }()
	var opt Option
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Println("rpc server: options error: ", err)
		return
	}
	if opt.MagicNumber != MagicNumber {
		log.Printf("rpc server: invalid magic number %x", opt.MagicNumber)
		return
	}
	f := codec.NewCodecFuncMap[opt.CodecType]
	server.serveCodec(f(conn))
}

// invalidRequest is a placeholder for response argv when error occurs
var invaldRequest = struct {
}{}

// serveCodec
/*
Three States
读取请求： readRequest
处理请求: handleRequest
回复请求：sendRequest
*/
func (server *Server) serveCodec(cc codec.Codec) {
	sending := new(sync.Mutex) // make sure to send a complete response
	wg := new(sync.WaitGroup)  // wait until all request are handled
	for {
		req, err := server.readRequest(cc)
		// 错误处理
		if err != nil {
			if req == nil {
				break // it`s not possible to recover, so close the connection
			}

			req.h.Error = err.Error()
			server.sendResponse(cc, req.h, invaldRequest, sending)
			continue
		}
		wg.Add(1)
		go server.handleRequest(cc, req, sending, wg)
	}

	wg.Wait()
	_ = cc.Close()
}

// request stores all information of a call
type request struct {
	h            *codec.Header // header of request
	argv, replyv reflect.Value // argv and replyv of request
}

func (server *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server: read error:", err)
		}
		return nil, err
	}
	return &h, nil
}

func (server *Server) readRequest(cc codec.Codec) (*request, error) {
	h, err := server.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}

	req := &request{h: h}
	// TODO: now we don`t know the type of request argv
	// day 1,just support it`s string
	req.argv = reflect.New(reflect.TypeOf(""))
	if err = cc.ReadBody(req.argv.Interface()); err != nil {
		log.Println("rpc server: read argv err:", err)
	}
	return req, nil
}

func (server *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		log.Println("rpc server: write response error:", err)
	}
}

func (server *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	// TODO, should call registered rpc methods to get the right replyv
	// day 1, just print argv and send a hello message
	defer wg.Done()
	log.Println(req.h, req.argv.Elem())
	//  body 作为字符串处理。接收到请求，打印 header，并回复 geerpc resp ${req.h.Seq}
	req.replyv = reflect.ValueOf(fmt.Sprintf("neerpc reps %d", req.h.Seq))

	server.sendResponse(cc, req.h, req.replyv.Interface(), sending)
}

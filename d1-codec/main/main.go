package main

import (
	"encoding/json"
	"fmt"
	"log"
	neerpc "neerpc"
	"neerpc/codec"
	"net"
	"time"
)

func startServer(addr chan string) {
	// pick a free port
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network, error:", err)
	}
	log.Println("Start rpc server on", l.Addr())
	addr <- l.Addr().String()
	neerpc.Accept(l)
}

func main() {
	addr := make(chan string)
	go startServer(addr)

	// in a fact, following code is like a simple neerpc client
	conn, _ := net.Dial("tcp", <-addr)
	defer func() { _ = conn.Close() }()

	time.Sleep(time.Second)
	// send options
	_ = json.NewEncoder(conn).Encode(neerpc.DdefaultOption)
	cc := codec.NewGobCodec(conn)
	// send request & receive response
	for i := 0; i < 5; i++ {
		h := &codec.Header{
			ServiceMethod: "Foo.Sum",
			Seq:           uint64(i),
		}
		_ = cc.Write(h, fmt.Sprintf("neerpc req %d", h.Seq))
		_ = cc.ReadHeader(h)
		var reply string
		_ = cc.ReadBody(&reply)
		log.Println("rpely:", reply)

	}

}

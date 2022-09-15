package main

import (
	"fmt"
	"log"
	neerpc "neerpc"
	"net"
	"sync"
	"time"
)

type Foo int

type Args struct {
	Num1, Num2 int
}

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func startServer(addr chan string) {
	var foo Foo
	if err := neerpc.Register(&foo); err != nil {
		log.Fatal("register error:", err)
	}

	// pick a free port
	l, err := net.Listen("tcp", ":6677")
	if err != nil {
		log.Fatal("network error:", err)
	}

	log.Println("start rpc server on", l.Addr().String())
	addr <- l.Addr().String()
	neerpc.Accept(l)
}

func main() {
	//log.SetFlags(0)
	//addr := make(chan string)
	//go startServer(addr)
	//<-addr
	//for {
	//	time.Sleep(time.Second)
	//}
	//client, _ := neerpc.Dial("tcp", <-addr)
	//defer func() { _ = client.Close() }()
	//
	//time.Sleep(time.Second)
	//// send request & receive response
	//var wg sync.WaitGroup
	//
	//for i := 0; i < 500; i++ {
	//	wg.Add(1)
	//	go func(i int) {
	//		defer wg.Done()
	//		args := &Args{Num1: i, Num2: i * i}
	//		var reply int
	//		if err := client.Call("Foo.Sum", args, &reply); err != nil {
	//			log.Fatal("call Foo.Sum error:", err)
	//		}
	//
	//		log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
	//	}(i)
	//}
	//wg.Wait()

	client, err := neerpc.Dial("tcp", "82.157.193.4:6677")
	if err != nil {
		fmt.Println(err)
	}
	defer func() { _ = client.Close() }()

	time.Sleep(time.Second)
	// send request & receive response
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := &Args{Num1: i, Num2: i * i}
			var reply int
			if err := client.Call("Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}

			log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
		}(i)
	}
	wg.Wait()
}

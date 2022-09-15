package main

import (
	"context"
	"log"
	"math"
	neerpc "neerpc"
	"net"
	"net/http"
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

func (f *Foo) Pow(args Args, reply *int) error {
	*reply = int(math.Pow10(args.Num1))
	return nil
}

func startServer(addrCh chan string) {
	var foo Foo
	l, _ := net.Listen("tcp", ":9999")
	_ = neerpc.Register(&foo)
	neerpc.HandleHTTP()
	addrCh <- l.Addr().String()
	_ = http.Serve(l, nil)
}

func call(addrCh chan string) {
	client, _ := neerpc.DialHTTP("tcp", <-addrCh)
	defer func() { _ = client.Close() }()

	time.Sleep(time.Second)
	// send request & receive response
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := &Args{Num1: i, Num2: i * i}
			var reply int
			if err := client.Call(context.Background(), "Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Printf("Foo.Sum: %d + %d = %d", args.Num1, args.Num2, reply)
			if err := client.Call(context.Background(), "Foo.Pow", args, &reply); err != nil {
				log.Fatal("call Foo.Pow error:", err)
			}
			log.Printf("Foo.Pow: 10^%d = %d", args.Num1, reply)
		}(i)
	}
	wg.Wait()
}
func main() {
	log.SetFlags(0)
	ch := make(chan string)
	go call(ch)
	startServer(ch)
}

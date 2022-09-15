package main

import (
	"log"
	neerpc "neerpc"
	"sync"
	"testing"
	"time"
)

func BenchmarkName(b *testing.B) {
	//log.SetFlags(0)
	//addr := make(chan string)
	//go startServer(addr)
	client, _ := neerpc.Dial("tcp", "82.157.193.4:6677")
	defer func() { _ = client.Close() }()

	time.Sleep(time.Second)
	// send request & receive response
	var wg sync.WaitGroup

	//fmt.Println("start")
	//b.ResetTimer()
	//b.StartTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < 500; i++ {
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
	//b.StopTimer()

}

package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

type P struct {
	X, Y, Z int
	Name    string
}

type Q struct {
	X, Y *int32
	Name string
}

// 此实例显示了包的基本用法：创建编码器
// 传输一些值，用解码器接收
func main() {
	// 初始化编码和解码器，通常是enc和dec
	// 绑定到网络连接和编码器和解码器会在不同的进程中运行
	var network bytes.Buffer // 替代网络连接
	enc := gob.NewEncoder(&network)
	dec := gob.NewDecoder(&network)
	// Encoding （发送）一些值
	err := enc.Encode(P{3, 4, 5, "Pythagoras"})
	if err != nil {
		log.Fatal("encode error:", err)
	}
	err = enc.Encode(P{1782, 1841, 1922, "Treehouse"})
	if err != nil {
		log.Fatal("encode error:", err)
	}

	// Decode (接收) 并打印值
	var p P
	err = dec.Decode(&p)
	//fmt.Printf("%q: {%d, %d, %d}\n", p.Name, p.X, p.Y, p.Z)
	var q Q
	err = dec.Decode(&q)
	if err != nil {
		log.Fatal("decode error 1:", err)
	}

	fmt.Printf("%q: {%d, %d}\n", q.Name, *q.X, *q.Y)
	err = dec.Decode(&q)
	if err != nil {
		log.Fatal("decode error 2:", err)
	}
	fmt.Printf("%q: {%d, %d}\n", q.Name, *q.X, *q.Y)

}

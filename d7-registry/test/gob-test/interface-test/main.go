package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"math"
)

type Point struct {
	X, Y int
}

func (p Point) Hypotenuse() float64 {
	return math.Hypot(float64(p.X), float64(p.Y))
}

type Pythagoras interface {
	Hypotenuse() float64
}

// 此实例显示如何编码接口值，关键的与常规类型的区别是注册具体类型，实现接口
func main() {
	var network bytes.Buffer // 替代 (Stand-in) 网络

	// 必须注册编码器和解码器的具体类型（这将是通常与编码器不同的机器上）。
	// 在每一端，告诉了发送具体类型的引擎实现接口
	gob.Register(Point{})

	//  创建编码器并发送一些值
	enc := gob.NewEncoder(&network)
	for i := 0; i <= 3; i++ {
		interfaceEncode(enc, Point{3 * i, 4 * i})
	}

	// 创建解码器并接收一些值
	dec := gob.NewDecoder(&network)
	for i := 0; i <= 3; i++ {
		result := interfaceDecode(dec)
		fmt.Println(result.Hypotenuse())
	}
}

// interfaceEncode 将接口值编码到编码器中
func interfaceEncode(enc *gob.Encoder, p Pythagoras) {
	// 除非具体类型，否则编码将失败
	// 注册， 在调用函数中注册了
	// 将指针传递给接口，以便Encode看到（并因发送）
	// 一个值界面类型，如果直接传递p， 它会看到具体的类型
	// 有关背景，请参阅博客文章“（The Laws of Reflection）反思的法则”。
	err := enc.Encode(&p)
	if err != nil {
		log.Fatal("encode:", err)
	}
}

// interfaceDecode 解码流中的下一个接口值并返回
func interfaceDecode(dec *gob.Decoder) Pythagoras {
	// 除非线路上的具体类型已经解码， 否则解码器将失效
	// 注册， 在调用函数中注册了
	var p Pythagoras
	err := dec.Decode(&p)
	if err != nil {
		log.Fatal("decode:", err)
	}

	return p
}

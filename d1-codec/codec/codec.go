package codec

import "io"

/*
information encoding and decoding communication
信息编解码
*/

type Header struct {
	ServiceMethod string // format "service.Method"	 // 服务名和方法名
	Seq           uint64 // sequence number chose by client	// 请求的序号，某个请求的ID。区分不同的请求
	Error         string // 错误信息
}

type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

type NewCodecFunc func(closer io.ReadWriteCloser) Codec

type Type string

const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json" // not implemented
)

var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}

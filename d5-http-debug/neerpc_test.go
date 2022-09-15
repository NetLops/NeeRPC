package NeeRPC

import (
	"errors"
	"fmt"
	"neerpc/codec"
	"reflect"
	"testing"
)

func TestReflect(t *testing.T) {
	s := "离谱"
	revalue := reflect.ValueOf(&s)
	fmt.Println(revalue.Elem())
	fmt.Println(revalue.Interface())
}

func TestHeader(t *testing.T) {
	var header codec.Header
	header.Seq = 1
	header.Error = errors.New("test").Error()
	header.ServiceMethod = "test"
	fmt.Println(header)
}

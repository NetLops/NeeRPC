package d1_codec

import (
	"fmt"
	"reflect"
	"testing"
)

func TestReflect(t *testing.T) {
	s := "离谱"
	revalue := reflect.ValueOf(&s)
	fmt.Println(revalue.Elem())
	fmt.Println(revalue.Interface())
}

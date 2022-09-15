package reflect_test

import (
	"fmt"
	"reflect"
	"testing"
)

func TestReflectInterfaceElem(t *testing.T) {
	s := "test"
	//fmt.Println(reflect.ValueOf(s).Interface())
	//fmt.Println(reflect.ValueOf(&s).Elem().Interface())
	//fmt.Println(reflect.ValueOf(&s).Interface())

	fmt.Println(reflect.Indirect(reflect.ValueOf(s)).Type().Name())
	fmt.Println(reflect.Indirect(reflect.ValueOf(&s)))
	fmt.Println(reflect.Indirect(reflect.ValueOf(reflect.ValueOf(&s).Elem())))
}

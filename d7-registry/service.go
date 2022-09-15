package NeeRPC

import (
	"go/ast"
	"log"
	"reflect"
	"sync/atomic"
)

/*
the method’s type is exported. – 方法所属类型是导出的。
the method is exported. – 方式是导出的。
the method has two arguments, both exported (or builtin) types. – 两个入参，均为导出或内置类型。
the method’s second argument is a pointer. – 第二个入参必须是一个指针。
the method has return type error. – 返回值为 error 类型。
*/
// func (t *T) MethodName(argType T1, replyType *T2) error

type methodType struct {
	method    reflect.Method // 方法本身
	ArgType   reflect.Type   // 第一个参数的类型
	ReplyType reflect.Type   // 第二个参数的类型
	numCalls  uint64         // 统计方法调用次数
}

func (m *methodType) NumCalls() uint64 {
	return atomic.LoadUint64(&m.numCalls)
}

func (m *methodType) newArgv() reflect.Value {
	var argv reflect.Value
	// arg may be a pointer type, or a value type
	if m.ArgType.Kind() == reflect.Ptr {
		argv = reflect.New(m.ArgType.Elem())
	} else {
		argv = reflect.New(m.ArgType).Elem()
	}
	return argv
}

func (m *methodType) newReply() reflect.Value {
	// replyv must be a pointer type
	replyv := reflect.New(m.ReplyType.Elem())
	switch m.ReplyType.Elem().Kind() {
	case reflect.Map:
		replyv.Elem().Set(reflect.MakeMap(m.ReplyType.Elem()))
	case reflect.Slice:
		replyv.Elem().Set(reflect.MakeSlice(m.ReplyType.Elem(), 0, 0))
	}
	return replyv
}

type service struct {
	name   string                 // 映射结构体的名称
	typ    reflect.Type           //结构体的类型
	rcvr   reflect.Value          // 结构体的实例本身， 在调用时，需要rcvr 作为第0个参数
	method map[string]*methodType // 存储映射的结构体的所有符合条件的方法
}

func newService(rcvr any) *service {
	s := new(service)
	s.rcvr = reflect.ValueOf(rcvr)
	s.name = reflect.Indirect(s.rcvr).Type().Name()
	s.typ = reflect.TypeOf(rcvr)
	if !ast.IsExported(s.name) {
		log.Fatalf("rpc server: %s is not a valid service name", s.name)
	}
	s.registerMethods()
	return s
}

// registerMethods
// 两个导出或内置类型的入参（反射时为 3 个，第 0 个是自身，类似于 python 的 self，java 中的 this）
// 返回值有且只有 1 个，类型为 error
func (s *service) registerMethods() {
	s.method = make(map[string]*methodType)
	for i := 0; i < s.typ.NumMethod(); i++ {
		method := s.typ.Method(i)
		mType := method.Type
		// 要符合这个
		// func (t *T) MethodName(argType T1, replyType *T2) error
		if mType.NumIn() != 3 || mType.NumOut() != 1 {
			continue
		}
		// 0 就不不需要了 毕竟那是 s.rcvr 自己
		argType, replyType := mType.In(1), mType.In(2)
		if !isExportOrBuiltinType(argType) || !isExportOrBuiltinType(replyType) {
			continue
		}
		s.method[method.Name] = &methodType{
			method:    method,
			ArgType:   argType,
			ReplyType: replyType,
		}
		log.Printf("rpc server: register %s.%s\n", s.name, method.Name)
	}
}

func isExportOrBuiltinType(t reflect.Type) bool {
	return ast.IsExported(t.Name()) || t.PkgPath() == "" // "" 当前 在自己包里的， 自家人就不讲究是否对外公开了
}

func (s *service) call(m *methodType, argv, replyv reflect.Value) error {
	atomic.AddUint64(&m.numCalls, 1)
	f := m.method.Func
	returnValues := f.Call([]reflect.Value{s.rcvr, argv, replyv})
	if errInter := returnValues[0].Interface(); errInter != nil {
		return errInter.(error)
	}
	return nil
}

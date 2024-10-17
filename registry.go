package minrpc

import (
    "fmt"
    "reflect"
    "sync"
)

// RPCMethod 接口表示一个RPC方法
type RPCMethod interface {
    Name() string
    Call(args []reflect.Value) (reflect.Value, error)
}

// Registry 是一个注册中心，负责管理所有注册的RPC方法
type Registry struct {
    methods map[string]RPCMethod
    mu      sync.RWMutex
}

// NewRegistry 创建一个新的注册中心
func NewRegistry() *Registry {
    return &Registry{
        methods: make(map[string]RPCMethod),
    }
}

// RegisterMethod 注册一个RPC方法
func (r *Registry) RegisterMethod(method RPCMethod) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.methods[method.Name()] = method
}

// GetMethod 获取一个RPC方法
func (r *Registry) GetMethod(name string) (RPCMethod, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    method, ok := r.methods[name]
    return method, ok
}

// wrapFunction 将普通函数包装成符合RPCMethod接口的结构体
func WrapFunction(name string, fn interface{}) RPCMethod {
    fnValue := reflect.ValueOf(fn)
    fnType := fnValue.Type()

    return &functionWrapper{
        name: name,
        fn:   fnValue,
        fnType: fnType,
    }
}

// functionWrapper 是一个包装了普通函数的结构体
type functionWrapper struct {
    name   string
    fn     reflect.Value
    fnType reflect.Type
}

func (w *functionWrapper) Name() string {
    return w.name
}

func (w *functionWrapper) Call(args []reflect.Value) (reflect.Value, error) {
    if len(args) != w.fnType.NumIn() {
        return reflect.Value{}, fmt.Errorf("expected %d arguments, got %d", w.fnType.NumIn(), len(args))
    }

    result := w.fn.Call(args)
    if len(result) == 0 {
        return reflect.Value{}, nil
    }

    if len(result) == 1 {
        return result[0], nil
    }

    if len(result) == 2 {
        if result[1].IsNil() {
            return result[0], nil
        }
        return result[0], result[1].Interface().(error)
    }

    return reflect.Value{}, fmt.Errorf("unexpected number of return values")
}
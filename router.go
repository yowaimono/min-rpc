package minrpc

import (
    "encoding/json"
    "fmt"
    "net/http"
    "reflect"
)

// RPCServer 是一个简单的RPC服务器
type RPCServer struct {
    registry *Registry
}

// NewRPCServer 创建一个新的RPC服务器
func NewRPCServer(registry *Registry) *RPCServer {
    return &RPCServer{
        registry: registry,
    }
}

// ServeHTTP 处理HTTP请求
func (s *RPCServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    var req Request
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    method, ok := s.registry.GetMethod(req.Header.Method)
    if !ok {
        http.Error(w, fmt.Sprintf("method %s not found", req.Header.Method), http.StatusNotFound)
        return
    }

    // 将请求体转换为 reflect.Value 类型
    var args []interface{}
    if err := json.Unmarshal(req.Body, &args); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    reflectArgs := make([]reflect.Value, len(args))
    for i, arg := range args {
        reflectArgs[i] = reflect.ValueOf(arg)
    }

    // 调用方法
    result, err := method.Call(reflectArgs)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // 将结果编码为响应体
    respBody, err := json.Marshal(result.Interface())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // 编码响应
    resp, err := EncodeResponse(req.Header.RequestID, "success", respBody)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(resp)
}
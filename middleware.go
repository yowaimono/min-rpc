package minrpc

import (
    "log"
    "net/http"
    "time"
)

// Middleware 是一个中间件接口
type Middleware interface {
    Handle(next http.Handler) http.Handler
}

// MiddlewareChain 是一个中间件链
type MiddlewareChain struct {
    middlewares []Middleware
}

// NewMiddlewareChain 创建一个新的中间件链
func NewMiddlewareChain() *MiddlewareChain {
    return &MiddlewareChain{
        middlewares: make([]Middleware, 0),
    }
}

// Use 注册一个中间件
func (c *MiddlewareChain) Use(m Middleware) *MiddlewareChain {
    c.middlewares = append(c.middlewares, m)
    return c
}

// Then 将中间件链应用到处理器
func (c *MiddlewareChain) Then(h http.Handler) http.Handler {
    for i := len(c.middlewares) - 1; i >= 0; i-- {
        h = c.middlewares[i].Handle(h)
    }
    return h
}

// LoggingMiddleware 是一个日志中间件
type LoggingMiddleware struct{}

func (m *LoggingMiddleware) Handle(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
    })
}

// AuthMiddleware 是一个认证中间件
type AuthMiddleware struct{}

func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 在这里实现认证逻辑
        next.ServeHTTP(w, r)
    })
}
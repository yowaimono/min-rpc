package minrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// HTTPClientPool 是一个HTTP客户端连接池
type HTTPClientPool struct {
	clients chan *http.Client
}

// NewHTTPClientPool 创建一个新的HTTP客户端连接池
func NewHTTPClientPool(size int) *HTTPClientPool {
	pool := &HTTPClientPool{
		clients: make(chan *http.Client, size),
	}

	for i := 0; i < size; i++ {
		pool.clients <- &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	return pool
}

// Get 从连接池中获取一个HTTP客户端
func (p *HTTPClientPool) Get() *http.Client {
	return <-p.clients
}

// Put 将HTTP客户端放回连接池
func (p *HTTPClientPool) Put(client *http.Client) {
	p.clients <- client
}

// RPCClient 是一个简单的RPC客户端
type RPCClient struct {
	endpoint   string
	pool       *HTTPClientPool
	serializer Serializer
}

// NewRPCClient 创建一个新的RPC客户端
func NewRPCClient(endpoint string, pool *HTTPClientPool, serializer Serializer) *RPCClient {
	return &RPCClient{
		endpoint:   endpoint,
		pool:       pool,
		serializer: serializer,
	}
}

// Call 发送RPC请求并返回结果
func (c *RPCClient) Call(method string, args ...interface{}) (interface{}, error) {
	requestID := fmt.Sprintf("%d", time.Now().UnixNano())

	// 编码请求
	reqBody, err := EncodeRequest(requestID, method, args)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(reqBody)

	client := c.pool.Get()
	defer c.pool.Put(client)

	resp, err := client.Post(c.endpoint, "application/json", buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var respData Response
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, err
	}

	if respData.Header.Status != "success" {
		return nil, fmt.Errorf("RPC error: %s", respData.Header.Error)
	}

	var result interface{}
	if err := json.Unmarshal(respData.Body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

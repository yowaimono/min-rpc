package minrpc

import (
	"encoding/json"
)

// RequestHeader 表示请求头
type RequestHeader struct {
	RequestID string `json:"request_id"`
	Method    string `json:"method"`
}

// ResponseHeader 表示响应头
type ResponseHeader struct {
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
	Error     string `json:"error"`
}

// Request 表示一个RPC请求
type Request struct {
	Header RequestHeader   `json:"header"`
	Body   json.RawMessage `json:"body"`
}

// Response 表示一个RPC响应
type Response struct {
	Header ResponseHeader  `json:"header"`
	Body   json.RawMessage `json:"body"`
}

// EncodeRequest 编码请求
func EncodeRequest(requestID, method string, body interface{}) ([]byte, error) {
	bodyData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req := Request{
		Header: RequestHeader{
			RequestID: requestID,
			Method:    method,
		},
		Body: bodyData,
	}

	return json.Marshal(req)
}

// DecodeRequest 解码请求
func DecodeRequest(data []byte) (*Request, error) {
	var req Request
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

// EncodeResponse 编码响应
func EncodeResponse(requestID, status string, body interface{}) ([]byte, error) {
	bodyData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp := Response{
		Header: ResponseHeader{
			RequestID: requestID,
			Status:    status,
			Error:     "",
		},
		Body: bodyData,
	}

	return json.Marshal(resp)
}

// DecodeResponse 解码响应
func DecodeResponse(data []byte) (*Response, error) {
	var resp Response
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Serializer 是一个序列化接口
type Serializer interface {
	Serialize(interface{}) ([]byte, error)
	Deserialize([]byte, interface{}) error
}

// JSONSerializer 是一个JSON序列化实现
type JSONSerializer struct{}

func (s *JSONSerializer) Serialize(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (s *JSONSerializer) Deserialize(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

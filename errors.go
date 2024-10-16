package minrpc

import (
    "fmt"
    "log"
)

// ErrorCode 是一个错误码类型
type ErrorCode string

const (
    ErrMethodNotFound ErrorCode = "MethodNotFound"
    ErrInvalidRequest ErrorCode = "InvalidRequest"
    ErrInternalServer ErrorCode = "InternalServerError"
)

// Error 是一个错误结构体
type Error struct {
    Code    ErrorCode
    Message string
}

func (e *Error) Error() string {
    return fmt.Sprintf("Error: %s - %s", e.Code, e.Message)
}

// NewError 创建一个新的错误
func NewError(code ErrorCode, message string) *Error {
    return &Error{
        Code:    code,
        Message: message,
    }
}

// LogError 记录错误日志
func LogError(err error) {
    log.Printf("Error: %v", err)
}
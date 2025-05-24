package utils

import "errors"

// 定义错误类型
var (
	ErrInvalidParams     = errors.New("无效的参数")
	ErrUnauthorized      = errors.New("未授权")
	ErrForbidden         = errors.New("禁止访问")
	ErrNotFound          = errors.New("资源不存在")
	ErrInternalServer    = errors.New("服务器内部错误")
	ErrDatabaseOperation = errors.New("数据库操作失败")
	ErrRedisOperation    = errors.New("Redis操作失败")
	ErrOSSOperation      = errors.New("OSS操作失败")
	ErrAIOperation       = errors.New("AI操作失败")
	ErrTaskOperation     = errors.New("任务操作失败")
)

// ErrorCode 错误码
const (
	CodeSuccess          = 200
	CodeInvalidParams    = 400
	CodeUnauthorized     = 401
	CodeForbidden        = 403
	CodeNotFound         = 404
	CodeInternalServer   = 500
	CodeDatabaseError    = 501
	CodeRedisError       = 502
	CodeOSSError         = 503
	CodeAIError          = 504
	CodeTaskError        = 505
)

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewError 创建错误响应
func NewError(code int, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: message,
	}
}

// IsError 判断是否为错误响应
func IsError(err error) bool {
	return err != nil
}

// GetErrorCode 获取错误码
func GetErrorCode(err error) int {
	switch err {
	case ErrInvalidParams:
		return CodeInvalidParams
	case ErrUnauthorized:
		return CodeUnauthorized
	case ErrForbidden:
		return CodeForbidden
	case ErrNotFound:
		return CodeNotFound
	case ErrDatabaseOperation:
		return CodeDatabaseError
	case ErrRedisOperation:
		return CodeRedisError
	case ErrOSSOperation:
		return CodeOSSError
	case ErrAIOperation:
		return CodeAIError
	case ErrTaskOperation:
		return CodeTaskError
	default:
		return CodeInternalServer
	}
} 
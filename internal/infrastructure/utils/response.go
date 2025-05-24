package utils

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Response 统一响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(r *ghttp.Request, data interface{}) {
	r.Response.WriteJsonExit(Response{
		Code: 200,
		Msg:  "操作成功",
		Data: data,
	})
}

// Error 错误响应
func Error(r *ghttp.Request, code int, msg string) {
	r.Response.WriteJsonExit(Response{
		Code: code,
		Msg:  msg,
	})
}

// ValidationError 参数验证错误响应
func ValidationError(r *ghttp.Request, err error) {
	r.Response.WriteJsonExit(Response{
		Code: 400,
		Msg:  err.Error(),
	})
}

// ServerError 服务器错误响应
func ServerError(r *ghttp.Request, err error) {
	g.Log().Error(r.GetCtx(), err)
	r.Response.WriteJsonExit(Response{
		Code: 500,
		Msg:  "服务器内部错误",
	})
}

// UnauthorizedError 未授权错误响应
func UnauthorizedError(r *ghttp.Request) {
	r.Response.WriteJsonExit(Response{
		Code: 401,
		Msg:  "未授权",
	})
}

// ForbiddenError 禁止访问错误响应
func ForbiddenError(r *ghttp.Request) {
	r.Response.WriteJsonExit(Response{
		Code: 403,
		Msg:  "禁止访问",
	})
}

// NotFoundError 资源不存在错误响应
func NotFoundError(r *ghttp.Request) {
	r.Response.WriteJsonExit(Response{
		Code: 404,
		Msg:  "资源不存在",
	})
} 
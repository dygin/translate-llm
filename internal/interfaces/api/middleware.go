package api

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/go-redis/redis/v8"
	"ai-translate/internal/utils"
	"context"
	"io/ioutil"
	"strings"
	"time"
)

// Auth 认证中间件
func Auth(r *ghttp.Request) {
	// 获取认证令牌
	token := r.Header.Get("Authorization")
	if token == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    utils.ErrUnauthorized,
			Message: "未授权访问",
		})
		r.Exit()
		return
	}

	// 验证令牌格式
	if !strings.HasPrefix(token, "Bearer ") {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    utils.ErrUnauthorized,
			Message: "无效的令牌格式",
		})
		r.Exit()
		return
	}

	// 提取令牌
	token = strings.TrimPrefix(token, "Bearer ")

	// 验证令牌
	claims, err := utils.ValidateToken(token)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    utils.ErrUnauthorized,
			Message: "令牌验证失败: " + err.Error(),
		})
		r.Exit()
		return
	}

	// 将用户信息存储到上下文
	r.SetCtxVar("user_id", claims.UserID)
	r.SetCtxVar("user_role", claims.Role)

	r.Middleware.Next()
}

// Logger 日志中间件
func Logger(r *ghttp.Request) {
	startTime := gtime.TimestampMilli()
	requestID := r.GetCtxVar("request_id").String()
	
	// 记录请求信息
	requestBody, _ := r.GetBody()
	requestHeaders := r.Header.Clone()
	requestHeaders.Del("Authorization") // 移除敏感信息
	
	logData := map[string]interface{}{
		"request_id":    requestID,
		"method":        r.Method,
		"path":          r.URL.Path,
		"query":         r.URL.RawQuery,
		"headers":       requestHeaders,
		"body":          string(requestBody),
		"client_ip":     r.GetClientIp(),
		"user_agent":    r.UserAgent(),
		"start_time":    startTime,
	}
	
	g.Log().Infof(r.Context(), "请求开始: %s", utils.JsonEncode(logData))
	
	// 创建响应写入器
	responseWriter := &responseBodyWriter{
		ResponseWriter: r.Response.Writer,
		body:          &bytes.Buffer{},
	}
	r.Response.Writer = responseWriter
	
	r.Middleware.Next()
	
	// 记录响应信息
	endTime := gtime.TimestampMilli()
	responseHeaders := r.Response.Header().Clone()
	
	logData = map[string]interface{}{
		"request_id":    requestID,
		"method":        r.Method,
		"path":          r.URL.Path,
		"status":        r.Response.Status,
		"headers":       responseHeaders,
		"body":          responseWriter.body.String(),
		"duration":      endTime - startTime,
		"end_time":      endTime,
	}
	
	g.Log().Infof(r.Context(), "请求结束: %s", utils.JsonEncode(logData))
}

// ErrorHandler 错误处理中间件
func ErrorHandler(r *ghttp.Request) {
	r.Middleware.Next()

	// 如果已经有响应，则不再处理
	if r.Response.BufferLength() > 0 {
		return
	}

	// 获取错误信息
	var (
		code    = r.Response.Status
		message = r.Response.BufferString()
	)

	// 根据状态码处理错误
	switch code {
	case 400:
		message = "请求参数错误"
	case 401:
		message = "未授权访问"
	case 403:
		message = "禁止访问"
	case 404:
		message = "请求的资源不存在"
	case 405:
		message = "请求方法不允许"
	case 429:
		message = "请求过于频繁"
	case 500:
		message = "服务器内部错误"
	case 502:
		message = "网关错误"
	case 503:
		message = "服务不可用"
	case 504:
		message = "网关超时"
	}

	// 返回错误响应
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// CORS 跨域中间件
func CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Response.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,X-Requested-With")
	r.Response.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
	r.Response.Header().Set("Access-Control-Allow-Credentials", "true")
	r.Response.Header().Set("Access-Control-Max-Age", "3600")
	
	r.Middleware.Next()
}

// RateLimit 限流中间件
func RateLimit(r *ghttp.Request) {
	// 获取Redis客户端
	redisClient := g.Redis().GetClient()
	
	// 获取客户端IP
	clientIP := r.GetClientIp()
	
	// 构建限流键
	key := fmt.Sprintf("rate_limit:%s:%s", clientIP, r.URL.Path)
	
	// 获取限流配置
	limit := g.Cfg().MustGet("rate_limit.limit").Int()     // 限制次数
	window := g.Cfg().MustGet("rate_limit.window").Int()   // 时间窗口(秒)
	
	// 使用Redis实现滑动窗口限流
	ctx := context.Background()
	now := time.Now().Unix()
	
	// 清理过期的请求记录
	redisClient.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", now-int64(window)))
	
	// 获取当前时间窗口内的请求数
	count, err := redisClient.ZCard(ctx, key).Result()
	if err != nil {
		g.Log().Errorf(r.Context(), "限流检查失败: %v", err)
		r.Middleware.Next()
		return
	}
	
	// 检查是否超过限制
	if count >= int64(limit) {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    429,
			Message: "请求过于频繁，请稍后再试",
		})
		r.Exit()
		return
	}
	
	// 记录本次请求
	redisClient.ZAdd(ctx, key, &redis.Z{
		Score:  float64(now),
		Member: now,
	})
	
	// 设置过期时间
	redisClient.Expire(ctx, key, time.Duration(window)*time.Second)
	
	r.Middleware.Next()
}

// RequestID 请求ID中间件
func RequestID(r *ghttp.Request) {
	// 生成请求ID
	requestID := utils.GenerateUUID()
	
	// 设置请求ID到上下文
	r.SetCtxVar("request_id", requestID)
	
	// 设置请求ID到响应头
	r.Response.Header().Set("X-Request-ID", requestID)
	
	r.Middleware.Next()
}

// Recovery 恢复中间件
func Recovery(r *ghttp.Request) {
	defer func() {
		if err := recover(); err != nil {
			// 记录错误日志
			g.Log().Errorf(r.Context(), "请求处理异常: %v\n堆栈信息: %s", err, utils.GetStack())
			
			// 返回错误响应
			r.Response.WriteJson(ghttp.DefaultHandlerResponse{
				Code:    utils.ErrInternalServer,
				Message: "服务器内部错误",
			})
		}
	}()
	
	r.Middleware.Next()
}

// ValidateRequest 请求参数验证中间件
func ValidateRequest(r *ghttp.Request) {
	// 获取请求参数
	var params map[string]interface{}
	if err := r.Parse(&params); err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "请求参数解析失败: " + err.Error(),
		})
		r.Exit()
		return
	}
	
	// 验证必填参数
	requiredFields := []string{"work_id", "batch_id", "type"}
	for _, field := range requiredFields {
		if _, ok := params[field]; !ok {
			r.Response.WriteJson(ghttp.DefaultHandlerResponse{
				Code:    400,
				Message: fmt.Sprintf("缺少必填参数: %s", field),
			})
			r.Exit()
			return
		}
	}
	
	// 验证参数类型
	if taskType, ok := params["type"].(string); ok {
		if !utils.IsValidTaskType(taskType) {
			r.Response.WriteJson(ghttp.DefaultHandlerResponse{
				Code:    400,
				Message: "无效的任务类型",
			})
			r.Exit()
			return
		}
	}
	
	// 验证优先级范围
	if priority, ok := params["priority"].(float64); ok {
		if priority < 0 || priority > 3 {
			r.Response.WriteJson(ghttp.DefaultHandlerResponse{
				Code:    400,
				Message: "优先级必须在0-3之间",
			})
			r.Exit()
			return
		}
	}
	
	r.Middleware.Next()
}

// Compress 响应压缩中间件
func Compress(r *ghttp.Request) {
	// 检查是否支持压缩
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		r.Middleware.Next()
		return
	}
	
	// 设置响应头
	r.Response.Header().Set("Content-Encoding", "gzip")
	r.Response.Header().Set("Vary", "Accept-Encoding")
	
	// 创建gzip写入器
	gzipWriter := gzip.NewWriter(r.Response.Writer)
	defer gzipWriter.Close()
	
	// 替换响应写入器
	responseWriter := &responseBodyWriter{
		ResponseWriter: r.Response.Writer,
		body:          &bytes.Buffer{},
	}
	r.Response.Writer = responseWriter
	
	r.Middleware.Next()
	
	// 压缩响应内容
	gzipWriter.Write(responseWriter.body.Bytes())
}

// responseBodyWriter 响应体写入器
type responseBodyWriter struct {
	ghttp.ResponseWriter
	body *bytes.Buffer
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// RegisterMiddleware 注册中间件
func RegisterMiddleware(s *ghttp.Server) {
	// 全局中间件
	s.Use(
		Recovery,     // 恢复中间件
		RequestID,    // 请求ID中间件
		Logger,       // 日志中间件
		CORS,         // 跨域中间件
		Compress,     // 响应压缩中间件
		ErrorHandler, // 错误处理中间件
	)

	// 需要认证的路由组中间件
	s.Group("/api/v1").Bind([]ghttp.GroupItem{
		{"ALL", "/*", Auth}, // 认证中间件
	})

	// 需要限流的路由组中间件
	s.Group("/api/v1").Bind([]ghttp.GroupItem{
		{"ALL", "/*", RateLimit}, // 限流中间件
	})

	// 需要参数验证的路由组中间件
	s.Group("/api/v1/tasks").Bind([]ghttp.GroupItem{
		{"POST", "/*", ValidateRequest}, // 参数验证中间件
	})
} 
package middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

// CORS 跨域中间件
func CORS(r *ghttp.Request) {
	// 设置允许的域名
	r.Response.CORSDefault()

	// 设置允许的请求方法
	r.Response.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	// 设置允许的请求头
	r.Response.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

	// 设置允许携带凭证
	r.Response.Header().Set("Access-Control-Allow-Credentials", "true")

	// 设置预检请求的缓存时间
	r.Response.Header().Set("Access-Control-Max-Age", "3600")

	// 处理预检请求
	if r.Method == "OPTIONS" {
		r.Response.WriteJsonExit(200, "ok")
	}

	r.Middleware.Next()
} 
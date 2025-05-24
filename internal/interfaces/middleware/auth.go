package middleware

import (
	"ai-translate/internal/application"
	"context"
	"errors"
	"github.com/gogf/gf/v2/crypto/gaes"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

// Auth 鉴权中间件
func Auth(r *ghttp.Request) {
	// 获取token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		r.Response.WriteJsonExit(g.Map{
			"code": 401,
			"msg":  "未授权",
		})
	}

	// 解析token
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		r.Response.WriteJsonExit(g.Map{
			"code": 401,
			"msg":  "无效的token格式",
		})
	}

	// 验证token
	token, err := jwt.Parse(tokenParts[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名方法")
		}
		return []byte(g.Cfg().MustGet(context.Background(), "jwt.secret").String()), nil
	})

	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 401,
			"msg":  "token验证失败",
		})
	}

	// 获取用户信息
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		r.Response.WriteJsonExit(g.Map{
			"code": 401,
			"msg":  "无效的token",
		})
	}

	// 将用户信息存储到上下文
	r.SetCtxVar("user_id", uint64(claims["user_id"].(float64)))
	r.SetCtxVar("username", claims["username"].(string))

	r.Middleware.Next()
}

// Encrypt 加密中间件
func Encrypt(r *ghttp.Request) {
	// 获取加密密钥
	key := []byte(g.Cfg().MustGet(context.Background(), "aes.key").String())

	// 加密请求体
	if r.Body.Len() > 0 {
		encrypted, err := gaes.Encrypt(r.Body.Bytes(), key)
		if err != nil {
			r.Response.WriteJsonExit(g.Map{
				"code": 500,
				"msg":  "加密失败",
			})
		}
		r.Body = ghttp.NewRequestBody(encrypted)
	}

	r.Middleware.Next()

	// 加密响应体
	if r.Response.BufferLength() > 0 {
		encrypted, err := gaes.Encrypt(r.Response.Buffer(), key)
		if err != nil {
			r.Response.WriteJsonExit(g.Map{
				"code": 500,
				"msg":  "加密失败",
			})
		}
		r.Response.ClearBuffer()
		r.Response.Write(encrypted)
	}
}

// Decrypt 解密中间件
func Decrypt(r *ghttp.Request) {
	// 获取加密密钥
	key := []byte(g.Cfg().MustGet(context.Background(), "aes.key").String())

	// 解密请求体
	if r.Body.Len() > 0 {
		decrypted, err := gaes.Decrypt(r.Body.Bytes(), key)
		if err != nil {
			r.Response.WriteJsonExit(g.Map{
				"code": 500,
				"msg":  "解密失败",
			})
		}
		r.Body = ghttp.NewRequestBody(decrypted)
	}

	r.Middleware.Next()

	// 解密响应体
	if r.Response.BufferLength() > 0 {
		decrypted, err := gaes.Decrypt(r.Response.Buffer(), key)
		if err != nil {
			r.Response.WriteJsonExit(g.Map{
				"code": 500,
				"msg":  "解密失败",
			})
		}
		r.Response.ClearBuffer()
		r.Response.Write(decrypted)
	}
} 
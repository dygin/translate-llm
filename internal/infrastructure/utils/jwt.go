package utils

import (
	"errors"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// Claims 自定义JWT声明
type Claims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	RoleType int    `json:"role_type"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
func GenerateToken(userID uint64, username string, roleType int) (string, error) {
	// 获取JWT配置
	secret := g.Cfg().MustGet("jwt.secret").String()
	expire := g.Cfg().MustGet("jwt.expire").Int()

	// 创建声明
	claims := Claims{
		UserID:   userID,
		Username: username,
		RoleType: roleType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expire) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// 生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken 解析JWT token
func ParseToken(tokenString string) (*Claims, error) {
	// 获取JWT配置
	secret := g.Cfg().MustGet("jwt.secret").String()

	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证token
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的token")
}

// ValidateToken 验证token是否有效
func ValidateToken(tokenString string) bool {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return false
	}

	// 检查token是否过期
	if claims.ExpiresAt.Before(time.Now()) {
		return false
	}

	return true
} 
package application

import (
	"ai-translate/internal/domain/user"
	"ai-translate/internal/infrastructure/persistence"
	"context"
	"errors"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type userService struct {
	userRepo user.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService() user.UserService {
	return &userService{
		userRepo: persistence.NewUserRepository(),
	}
}

func (s *userService) Register(username, password, email string) (*user.User, error) {
	// 检查用户名是否已存在
	existingUser, err := s.userRepo.FindByUsername(username)
	if err == nil && existingUser != nil {
		return nil, errors.New("用户名已存在")
	}

	// 加密密码
	encryptedPassword, err := gmd5.EncryptString(password)
	if err != nil {
		return nil, err
	}

	// 创建用户
	newUser := &user.User{
		Username:  username,
		Password:  encryptedPassword,
		Email:     email,
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 保存用户
	err = s.userRepo.Save(newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *userService) Login(username, password string) (string, error) {
	// 查找用户
	u, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", err
	}

	// 验证密码
	encryptedPassword, err := gmd5.EncryptString(password)
	if err != nil {
		return "", err
	}
	if encryptedPassword != u.Password {
		return "", errors.New("密码错误")
	}

	// 生成JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  u.ID,
		"username": u.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	// 签名token
	tokenString, err := token.SignedString([]byte(g.Cfg().MustGet(context.Background(), "jwt.secret").String()))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *userService) GetUserInfo(id uint64) (*user.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *userService) UpdateUser(user *user.User) error {
	return s.userRepo.Update(user)
}

func (s *userService) DeleteUser(id uint64) error {
	return s.userRepo.Delete(id)
} 
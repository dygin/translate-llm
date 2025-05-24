package user

import (
	"time"
)

// User 用户实体
type User struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Email     string    `json:"email"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Roles     []Role    `json:"roles"`
}

// Role 角色实体
type Role struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserRepository 用户仓储接口
type UserRepository interface {
	FindByID(id uint64) (*User, error)
	FindByUsername(username string) (*User, error)
	Save(user *User) error
	Update(user *User) error
	Delete(id uint64) error
}

// UserService 用户服务接口
type UserService interface {
	Register(username, password, email string) (*User, error)
	Login(username, password string) (string, error)
	GetUserInfo(id uint64) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id uint64) error
} 
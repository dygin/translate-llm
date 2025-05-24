package persistence

import (
	"context"
	"ai-translate/internal/domain/user"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type userRepository struct {
	db gdb.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository() user.UserRepository {
	return &userRepository{
		db: g.DB(),
	}
}

func (r *userRepository) FindByID(id uint64) (*user.User, error) {
	var u user.User
	err := r.db.Model("users").Where("id", id).Scan(&u)
	if err != nil {
		return nil, err
	}
	
	// 获取用户角色
	var roles []user.Role
	err = r.db.Model("user_roles ur").
		LeftJoin("roles r", "ur.role_id = r.id").
		Where("ur.user_id", id).
		Scan(&roles)
	if err != nil {
		return nil, err
	}
	u.Roles = roles
	
	return &u, nil
}

func (r *userRepository) FindByUsername(username string) (*user.User, error) {
	var u user.User
	err := r.db.Model("users").Where("username", username).Scan(&u)
	if err != nil {
		return nil, err
	}
	
	// 获取用户角色
	var roles []user.Role
	err = r.db.Model("user_roles ur").
		LeftJoin("roles r", "ur.role_id = r.id").
		Where("ur.user_id", u.ID).
		Scan(&roles)
	if err != nil {
		return nil, err
	}
	u.Roles = roles
	
	return &u, nil
}

func (r *userRepository) Save(user *user.User) error {
	_, err := r.db.Model("users").Insert(user)
	return err
}

func (r *userRepository) Update(user *user.User) error {
	_, err := r.db.Model("users").Where("id", user.ID).Update(user)
	return err
}

func (r *userRepository) Delete(id uint64) error {
	_, err := r.db.Model("users").Where("id", id).Delete()
	return err
} 
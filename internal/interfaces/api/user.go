package api

import (
	"ai-translate/internal/application"
	"ai-translate/internal/domain/user"
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type UserController struct {
	userService user.UserService
}

// NewUserController 创建用户控制器实例
func NewUserController() *UserController {
	return &UserController{
		userService: application.NewUserService(),
	}
}

// Register 用户注册
func (c *UserController) Register(r *ghttp.Request) {
	var req struct {
		Username string `json:"username" v:"required"`
		Password string `json:"password" v:"required"`
		Email    string `json:"email" v:"required|email"`
	}

	if err := r.Parse(&req); err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 400,
			"msg":  err.Error(),
		})
	}

	u, err := c.userService.Register(req.Username, req.Password, req.Email)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}

	r.Response.WriteJsonExit(g.Map{
		"code": 200,
		"msg":  "注册成功",
		"data": u,
	})
}

// Login 用户登录
func (c *UserController) Login(r *ghttp.Request) {
	var req struct {
		Username string `json:"username" v:"required"`
		Password string `json:"password" v:"required"`
	}

	if err := r.Parse(&req); err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 400,
			"msg":  err.Error(),
		})
	}

	token, err := c.userService.Login(req.Username, req.Password)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}

	r.Response.WriteJsonExit(g.Map{
		"code": 200,
		"msg":  "登录成功",
		"data": g.Map{
			"token": token,
		},
	})
}

// GetUserInfo 获取用户信息
func (c *UserController) GetUserInfo(r *ghttp.Request) {
	userID := r.GetCtxVar("user_id").Uint64()
	u, err := c.userService.GetUserInfo(userID)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}

	r.Response.WriteJsonExit(g.Map{
		"code": 200,
		"msg":  "获取成功",
		"data": u,
	})
}

// UpdateUser 更新用户信息
func (c *UserController) UpdateUser(r *ghttp.Request) {
	var req struct {
		Email string `json:"email" v:"email"`
	}

	if err := r.Parse(&req); err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 400,
			"msg":  err.Error(),
		})
	}

	userID := r.GetCtxVar("user_id").Uint64()
	u, err := c.userService.GetUserInfo(userID)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}

	if req.Email != "" {
		u.Email = req.Email
	}

	err = c.userService.UpdateUser(u)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}

	r.Response.WriteJsonExit(g.Map{
		"code": 200,
		"msg":  "更新成功",
		"data": u,
	})
}

// DeleteUser 删除用户
func (c *UserController) DeleteUser(r *ghttp.Request) {
	userID := r.GetCtxVar("user_id").Uint64()
	err := c.userService.DeleteUser(userID)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}

	r.Response.WriteJsonExit(g.Map{
		"code": 200,
		"msg":  "删除成功",
	})
} 
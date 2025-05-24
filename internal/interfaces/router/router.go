package router

import (
	"ai-translate/internal/interfaces/api"
	"ai-translate/internal/interfaces/middleware"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Register 注册路由
func Register(s *ghttp.Server) {
	// 全局中间件
	s.Use(
		middleware.Logger,    // 日志中间件
		middleware.CORS,      // 跨域中间件
		middleware.RateLimit, // 限流中间件
		middleware.Decrypt,   // 解密中间件
		middleware.Encrypt,   // 加密中间件
	)

	// 公开API
	s.Group("/api", func(group *ghttp.RouterGroup) {
		// 用户认证
		group.POST("/register", api.NewUserController().Register)
		group.POST("/login", api.NewUserController().Login)
	})

	// 需要认证的API
	s.Group("/api", func(group *ghttp.RouterGroup) {
		// 认证中间件
		group.Middleware(middleware.Auth)

		// 用户管理
		group.GET("/user/info", api.NewUserController().GetUserInfo)
		group.PUT("/user/info", api.NewUserController().UpdateUser)
		group.DELETE("/user", api.NewUserController().DeleteUser)

		// 作品管理
		group.POST("/works", api.NewWorkController().CreateWork)
		group.GET("/works/:id", api.NewWorkController().GetWork)
		group.GET("/works", api.NewWorkController().GetUserWorks)
		group.PUT("/works/:id", api.NewWorkController().UpdateWork)
		group.DELETE("/works/:id", api.NewWorkController().DeleteWork)

		// 翻译批次管理
		group.POST("/works/:id/batches", api.NewWorkController().CreateTranslationBatch)
		group.GET("/works/:id/batches/:batchId", api.NewWorkController().GetTranslationBatch)
		group.GET("/works/:id/batches", api.NewWorkController().GetWorkTranslationBatches)

		// 提示词管理
		group.POST("/prompts", api.NewPromptController().CreatePrompt)
		group.GET("/prompts/:id", api.NewPromptController().GetPrompt)
		group.GET("/prompts", api.NewPromptController().GetPromptsByType)
		group.PUT("/prompts/:id", api.NewPromptController().UpdatePrompt)
		group.DELETE("/prompts/:id", api.NewPromptController().DeletePrompt)

		// 任务管理
		group.POST("/tasks", api.NewTaskController().CreateTask)
		group.GET("/tasks/:id", api.NewTaskController().GetTask)
		group.GET("/tasks", api.NewTaskController().GetTasksByType)
		group.PUT("/tasks/:id", api.NewTaskController().UpdateTask)
		group.DELETE("/tasks/:id", api.NewTaskController().DeleteTask)
		group.POST("/tasks/:id/pause", api.NewTaskController().PauseTask)
		group.POST("/tasks/:id/resume", api.NewTaskController().ResumeTask)
		group.POST("/tasks/:id/retry", api.NewTaskController().RetryTask)
	})
} 
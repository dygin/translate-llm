package api

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

// RegisterRoutes 注册路由
func RegisterRoutes(s *ghttp.Server) {
	// 注册中间件
	RegisterMiddleware(s)

	// 创建任务控制器
	taskController, err := NewTaskController()
	if err != nil {
		panic(err)
	}

	// 任务管理路由组
	taskGroup := s.Group("/api/v1/tasks")
	{
		// 基础任务管理
		taskGroup.POST("/", taskController.Create)                    // 创建任务
		taskGroup.GET("/:id", taskController.Get)                     // 获取任务
		taskGroup.GET("/", taskController.List)                       // 获取任务列表
		taskGroup.PUT("/:id/status", taskController.UpdateStatus)     // 更新任务状态
		taskGroup.DELETE("/:id", taskController.Delete)               // 删除任务
		taskGroup.POST("/:id/retry", taskController.Retry)            // 重试任务
		taskGroup.GET("/stats", taskController.GetStats)              // 获取任务统计

		// 任务优先级管理
		taskGroup.PUT("/:id/priority", taskController.UpdatePriority) // 更新任务优先级
		taskGroup.PUT("/priority/batch", taskController.BatchUpdatePriority) // 批量更新任务优先级
		taskGroup.PUT("/priority/condition", taskController.UpdatePriorityByCondition) // 根据条件更新任务优先级

		// 优先级规则管理
		taskGroup.POST("/rules", taskController.CreatePriorityRule)   // 创建优先级规则
		taskGroup.GET("/rules/:id", taskController.GetPriorityRule)   // 获取优先级规则
		taskGroup.PUT("/rules/:id", taskController.UpdatePriorityRule) // 更新优先级规则
		taskGroup.DELETE("/rules/:id", taskController.DeletePriorityRule) // 删除优先级规则
		taskGroup.GET("/rules", taskController.ListPriorityRules)     // 获取优先级规则列表
		taskGroup.GET("/:id/priority-logs", taskController.GetPriorityLogs) // 获取优先级调整日志

		// 规则模板管理
		taskGroup.POST("/templates", taskController.CreateRuleTemplate) // 创建规则模板
		taskGroup.GET("/templates/:id", taskController.GetRuleTemplate) // 获取规则模板
		taskGroup.PUT("/templates/:id", taskController.UpdateRuleTemplate) // 更新规则模板
		taskGroup.DELETE("/templates/:id", taskController.DeleteRuleTemplate) // 删除规则模板
		taskGroup.GET("/templates", taskController.ListRuleTemplates) // 获取规则模板列表

		// 规则组管理
		taskGroup.POST("/groups", taskController.CreateRuleGroup)     // 创建规则组
		taskGroup.GET("/groups/:id", taskController.GetRuleGroup)     // 获取规则组
		taskGroup.PUT("/groups/:id", taskController.UpdateRuleGroup)  // 更新规则组
		taskGroup.DELETE("/groups/:id", taskController.DeleteRuleGroup) // 删除规则组
		taskGroup.GET("/groups", taskController.ListRuleGroups)       // 获取规则组列表
		taskGroup.POST("/groups/:id/rules", taskController.AddRuleToGroup) // 添加规则到组
		taskGroup.DELETE("/groups/:id/rules/:rule_id", taskController.RemoveRuleFromGroup) // 从组中移除规则
		taskGroup.GET("/groups/:id/rules", taskController.GetGroupRules) // 获取组中的规则
	}
} 
package main

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"os"
	"os/signal"
	"syscall"
	"ai-translate/internal/interfaces/router"
	"ai-translate/internal/service"
)

func main() {
	ctx := gctx.New()

	// 创建任务服务
	taskService, err := service.NewTaskService()
	if err != nil {
		g.Log().Fatalf(ctx, "创建任务服务失败: %v", err)
	}

	// 启动任务服务
	if err := taskService.Start(ctx); err != nil {
		g.Log().Fatalf(ctx, "启动任务服务失败: %v", err)
	}

	// 注册路由
	s := g.Server()
	router.Register(s)

	// 启动HTTP服务
	go func() {
		if err := s.Start(); err != nil {
			g.Log().Fatalf(ctx, "启动HTTP服务失败: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 优雅关闭
	g.Log().Info(ctx, "正在关闭服务...")

	// 停止任务服务
	taskService.Stop()

	// 关闭HTTP服务
	if err := s.Shutdown(); err != nil {
		g.Log().Errorf(ctx, "关闭HTTP服务失败: %v", err)
	}

	g.Log().Info(ctx, "服务已关闭")
} 
package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"
	"github.com/gogf/gf/v2/frame/g"
	"ai-translate/internal/infrastructure/ai"
	"ai-translate/internal/infrastructure/queue"
	"ai-translate/internal/model"
)

// TaskScheduler 任务调度器
type TaskScheduler struct {
	queue      queue.Queue
	repository model.TaskRepository
	aiDrivers  map[ai.DriverType]ai.AIService
	workers    int
	maxRetries int
	stopCh     chan struct{}
	wg         sync.WaitGroup
}

// NewTaskScheduler 创建任务调度器
func NewTaskScheduler(queue queue.Queue, repository model.TaskRepository, aiDrivers map[ai.DriverType]ai.AIService) (*TaskScheduler, error) {
	workers := g.Cfg().MustGet("queue.worker.numWorkers").Int()
	maxRetries := g.Cfg().MustGet("queue.worker.maxRetries").Int()

	return &TaskScheduler{
		queue:      queue,
		repository: repository,
		aiDrivers:  aiDrivers,
		workers:    workers,
		maxRetries: maxRetries,
		stopCh:     make(chan struct{}),
	}, nil
}

// Start 启动调度器
func (s *TaskScheduler) Start(ctx context.Context) error {
	// 启动工作协程
	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)
		go s.worker(ctx, i)
	}

	// 启动监控协程
	s.wg.Add(1)
	go s.monitor(ctx)

	return nil
}

// Stop 停止调度器
func (s *TaskScheduler) Stop() {
	close(s.stopCh)
	s.wg.Wait()
}

// worker 工作协程
func (s *TaskScheduler) worker(ctx context.Context, id int) {
	defer s.wg.Done()

	for {
		select {
		case <-s.stopCh:
			return
		default:
			// 获取待处理任务
			tasks, err := s.repository.GetPendingTasks(ctx, 1)
			if err != nil {
				g.Log().Errorf(ctx, "获取待处理任务失败: %v", err)
				time.Sleep(time.Second)
				continue
			}

			if len(tasks) == 0 {
				time.Sleep(time.Second)
				continue
			}

			task := tasks[0]

			// 更新任务状态为处理中
			if err := s.repository.UpdateStatus(ctx, task.ID, model.TaskStatusRunning); err != nil {
				g.Log().Errorf(ctx, "更新任务状态失败: %v", err)
				continue
			}

			// 处理任务
			if err := s.processTask(ctx, task); err != nil {
				g.Log().Errorf(ctx, "处理任务失败: %v", err)
				// 增加重试次数
				if err := s.repository.IncrementRetryCount(ctx, task.ID); err != nil {
					g.Log().Errorf(ctx, "增加重试次数失败: %v", err)
				}
				// 更新任务状态为失败
				if err := s.repository.UpdateStatus(ctx, task.ID, model.TaskStatusFailed); err != nil {
					g.Log().Errorf(ctx, "更新任务状态失败: %v", err)
				}
				continue
			}

			// 更新任务状态为已完成
			if err := s.repository.UpdateStatus(ctx, task.ID, model.TaskStatusCompleted); err != nil {
				g.Log().Errorf(ctx, "更新任务状态失败: %v", err)
			}
		}
	}
}

// processTask 处理任务
func (s *TaskScheduler) processTask(ctx context.Context, task *model.Task) error {
	// 获取AI驱动
	aiService, ok := s.aiDrivers[task.Driver]
	if !ok {
		return fmt.Errorf("不支持的AI驱动类型: %s", task.Driver)
	}

	var result string
	var err error

	// 根据任务类型处理
	switch task.Type {
	case model.TaskTypeContentGeneration:
		result, err = aiService.GenerateContent(ctx, task.Content, task.Language)
	case model.TaskTypeTranslation:
		result, err = aiService.Translate(ctx, task.Content, task.SourceLang, task.TargetLang)
	default:
		return fmt.Errorf("不支持的任务类型: %s", task.Type)
	}

	if err != nil {
		return err
	}

	// 更新任务结果
	task.Result = result
	task.CompletedAt = &time.Time{}
	return s.repository.Update(ctx, task)
}

// monitor 监控协程
func (s *TaskScheduler) monitor(ctx context.Context) {
	defer s.wg.Done()

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			// 获取任务统计
			stats, err := s.repository.GetStats(ctx, "")
			if err != nil {
				g.Log().Errorf(ctx, "获取任务统计失败: %v", err)
				continue
			}

			// 记录监控日志
			g.Log().Infof(ctx, "任务统计: 总数=%d, 待处理=%d, 处理中=%d, 已完成=%d, 失败=%d, 暂停=%d, 平均处理时间=%d秒, 成功率=%.2f%%",
				stats.TotalTasks,
				stats.PendingTasks,
				stats.RunningTasks,
				stats.CompletedTasks,
				stats.FailedTasks,
				stats.PausedTasks,
				stats.AvgProcessTime,
				stats.SuccessRate,
			)
		}
	}
} 
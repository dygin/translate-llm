package processor

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"sync"
	"time"
)

// TaskProcessor 任务处理器
type TaskProcessor struct {
	queue       Queue
	workers     int
	maxRetries  int
	wg          sync.WaitGroup
	stopChan    chan struct{}
	taskHandlers map[string]TaskHandler
}

// TaskHandler 任务处理函数类型
type TaskHandler func(ctx context.Context, taskID string, data []byte) error

// NewTaskProcessor 创建任务处理器
func NewTaskProcessor(queue Queue, workers, maxRetries int) *TaskProcessor {
	return &TaskProcessor{
		queue:       queue,
		workers:     workers,
		maxRetries:  maxRetries,
		stopChan:    make(chan struct{}),
		taskHandlers: make(map[string]TaskHandler),
	}
}

// RegisterHandler 注册任务处理函数
func (p *TaskProcessor) RegisterHandler(taskType string, handler TaskHandler) {
	p.taskHandlers[taskType] = handler
}

// Start 启动任务处理器
func (p *TaskProcessor) Start(ctx context.Context) error {
	// 启动工作协程
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker(ctx, i)
	}

	// 等待停止信号
	go func() {
		<-p.stopChan
		p.wg.Wait()
	}()

	return nil
}

// Stop 停止任务处理器
func (p *TaskProcessor) Stop() {
	close(p.stopChan)
}

// worker 工作协程
func (p *TaskProcessor) worker(ctx context.Context, id int) {
	defer p.wg.Done()

	g.Log().Infof(ctx, "Worker %d started", id)

	for {
		select {
		case <-ctx.Done():
			g.Log().Infof(ctx, "Worker %d stopped", id)
			return
		default:
			// 处理任务
			if err := p.processTask(ctx); err != nil {
				g.Log().Errorf(ctx, "Worker %d error: %v", id, err)
				time.Sleep(time.Second) // 错误后等待一秒
			}
		}
	}
}

// processTask 处理单个任务
func (p *TaskProcessor) processTask(ctx context.Context) error {
	// 从队列获取任务
	msg, err := p.queue.Consume(ctx)
	if err != nil {
		return fmt.Errorf("consume task failed: %v", err)
	}

	// 解析任务类型
	taskType := msg.Type
	handler, ok := p.taskHandlers[taskType]
	if !ok {
		return fmt.Errorf("unknown task type: %s", taskType)
	}

	// 处理任务
	err = handler(ctx, msg.ID, msg.Data)
	if err != nil {
		// 处理失败，重试
		if msg.Retries < p.maxRetries {
			msg.Retries++
			// 重新入队
			if err := p.queue.Publish(ctx, msg); err != nil {
				return fmt.Errorf("retry task failed: %v", err)
			}
		}
		return fmt.Errorf("process task failed: %v", err)
	}

	return nil
} 
package application

import (
	"ai-translate/internal/domain/task"
	"ai-translate/internal/infrastructure/persistence"
	"errors"
	"time"
)

type taskService struct {
	taskRepo  task.TaskRepository
	taskQueue task.TaskQueue
}

// NewTaskService 创建任务服务实例
func NewTaskService() task.TaskService {
	return &taskService{
		taskRepo:  persistence.NewTaskRepository(),
		taskQueue: persistence.NewTaskQueue(),
	}
}

func (s *taskService) CreateTask(task *task.Task) error {
	err := s.taskRepo.Save(task)
	if err != nil {
		return err
	}
	return s.taskQueue.Push(task)
}

func (s *taskService) GetTask(id uint64) (*task.Task, error) {
	return s.taskRepo.FindByID(id)
}

func (s *taskService) GetTasksByType(taskType int) ([]*task.Task, error) {
	return s.taskRepo.FindByType(taskType)
}

func (s *taskService) GetTasksByStatus(status int) ([]*task.Task, error) {
	return s.taskRepo.FindByStatus(status)
}

func (s *taskService) UpdateTask(task *task.Task) error {
	return s.taskRepo.Update(task)
}

func (s *taskService) DeleteTask(id uint64) error {
	// 从队列中移除任务
	err := s.taskQueue.Remove(id)
	if err != nil {
		return err
	}
	return s.taskRepo.Delete(id)
}

func (s *taskService) PauseTask(id uint64) error {
	t, err := s.taskRepo.FindByID(id)
	if err != nil {
		return err
	}

	// 从队列中移除任务
	err = s.taskQueue.Remove(id)
	if err != nil {
		return err
	}

	// 更新任务状态为暂停
	t.Status = 2 // 2:暂停
	t.UpdatedAt = time.Now()
	return s.taskRepo.Update(t)
}

func (s *taskService) ResumeTask(id uint64) error {
	t, err := s.taskRepo.FindByID(id)
	if err != nil {
		return err
	}

	// 更新任务状态为等待中
	t.Status = 0 // 0:等待中
	t.UpdatedAt = time.Now()
	err = s.taskRepo.Update(t)
	if err != nil {
		return err
	}

	// 重新加入队列
	return s.taskQueue.Push(t)
}

func (s *taskService) RetryTask(id uint64) error {
	t, err := s.taskRepo.FindByID(id)
	if err != nil {
		return err
	}

	// 检查重试次数
	if t.RetryCount >= t.MaxRetry {
		return errors.New("超过最大重试次数")
	}

	// 更新重试次数和状态
	t.RetryCount++
	t.Status = 0 // 0:等待中
	t.UpdatedAt = time.Now()
	err = s.taskRepo.Update(t)
	if err != nil {
		return err
	}

	// 重新加入队列
	return s.taskQueue.Push(t)
} 
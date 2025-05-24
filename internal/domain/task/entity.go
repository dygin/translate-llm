package task

import (
	"time"
)

// Task 任务实体
type Task struct {
	ID          uint64    `json:"id"`
	Type        int       `json:"type"` // 1:内容生成 2:翻译
	Priority    int       `json:"priority"`
	Status      int       `json:"status"`
	ReferenceID uint64    `json:"reference_id"` // 关联ID
	RetryCount  int       `json:"retry_count"`
	MaxRetry    int       `json:"max_retry"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TaskRepository 任务仓储接口
type TaskRepository interface {
	FindByID(id uint64) (*Task, error)
	FindByType(taskType int) ([]*Task, error)
	FindByStatus(status int) ([]*Task, error)
	Save(task *Task) error
	Update(task *Task) error
	Delete(id uint64) error
}

// TaskService 任务服务接口
type TaskService interface {
	CreateTask(task *Task) error
	GetTask(id uint64) (*Task, error)
	GetTasksByType(taskType int) ([]*Task, error)
	GetTasksByStatus(status int) ([]*Task, error)
	UpdateTask(task *Task) error
	DeleteTask(id uint64) error
	PauseTask(id uint64) error
	ResumeTask(id uint64) error
	RetryTask(id uint64) error
}

// TaskQueue 任务队列接口
type TaskQueue interface {
	Push(task *Task) error
	Pop() (*Task, error)
	Remove(id uint64) error
	GetLength() (int64, error)
	Clear() error
} 
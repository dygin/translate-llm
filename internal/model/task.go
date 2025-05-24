package model

import (
	"time"
	"ai-translate/internal/infrastructure/ai"
)

// TaskType 任务类型
type TaskType string

const (
	TaskTypeContentGeneration TaskType = "content_generation" // 内容生成
	TaskTypeTranslation      TaskType = "translation"        // 翻译
)

// TaskStatus 任务状态
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"   // 等待处理
	TaskStatusRunning   TaskStatus = "running"   // 处理中
	TaskStatusCompleted TaskStatus = "completed" // 已完成
	TaskStatusFailed    TaskStatus = "failed"    // 失败
	TaskStatusPaused    TaskStatus = "paused"    // 已暂停
)

// TaskPriority 任务优先级
type TaskPriority int

const (
	TaskPriorityLow    TaskPriority = 0 // 低优先级
	TaskPriorityNormal TaskPriority = 1 // 普通优先级
	TaskPriorityHigh   TaskPriority = 2 // 高优先级
	TaskPriorityUrgent TaskPriority = 3 // 紧急优先级
)

// Task 任务模型
type Task struct {
	ID          string       `json:"id" gorm:"primaryKey"`
	WorkID      string       `json:"work_id" gorm:"index"`
	BatchID     string       `json:"batch_id" gorm:"index"`
	Type        TaskType     `json:"type" gorm:"index"`
	Status      TaskStatus   `json:"status" gorm:"index"`
	Priority    TaskPriority `json:"priority" gorm:"index"` // 任务优先级
	Content     string       `json:"content"`
	Result      string       `json:"result"`
	Error       string       `json:"error"`
	Driver      ai.DriverType `json:"driver"`
	RetryCount  int          `json:"retry_count"`
	MaxRetries  int          `json:"max_retries"`
	StartedAt   *time.Time   `json:"started_at"`
	CompletedAt *time.Time   `json:"completed_at"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// TaskStats 任务统计
type TaskStats struct {
	TotalTasks      int64 `json:"total_tasks"`
	PendingTasks    int64 `json:"pending_tasks"`
	RunningTasks    int64 `json:"running_tasks"`
	CompletedTasks  int64 `json:"completed_tasks"`
	FailedTasks     int64 `json:"failed_tasks"`
	PausedTasks     int64 `json:"paused_tasks"`
	AvgProcessTime  int64 `json:"avg_process_time"` // 平均处理时间（秒）
	SuccessRate     float64 `json:"success_rate"`    // 成功率
}

// PriorityAdjustRule 优先级调整规则
type PriorityAdjustRule struct {
	ID          string       `json:"id" gorm:"primaryKey"`
	Name        string       `json:"name"`                    // 规则名称
	Description string       `json:"description"`             // 规则描述
	Condition   string       `json:"condition"`               // 触发条件
	Action      string       `json:"action"`                  // 执行动作
	Priority    TaskPriority `json:"priority"`                // 调整后的优先级
	WorkID      string       `json:"work_id" gorm:"index"`    // 工作ID
	BatchID     string       `json:"batch_id" gorm:"index"`   // 批次ID
	TaskType    TaskType     `json:"task_type" gorm:"index"`  // 任务类型
	Status      TaskStatus   `json:"status" gorm:"index"`     // 任务状态
	Enabled     bool         `json:"enabled"`                 // 是否启用
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// PriorityAdjustLog 优先级调整日志
type PriorityAdjustLog struct {
	ID        string       `json:"id" gorm:"primaryKey"`
	TaskID    string       `json:"task_id" gorm:"index"`
	RuleID    string       `json:"rule_id" gorm:"index"`
	OldPriority TaskPriority `json:"old_priority"`
	NewPriority TaskPriority `json:"new_priority"`
	Reason    string       `json:"reason"`
	CreatedAt time.Time    `json:"created_at"`
}

// TaskRepository 任务仓储接口
type TaskRepository interface {
	// 创建任务
	Create(ctx context.Context, task *Task) error
	// 获取任务
	Get(ctx context.Context, id string) (*Task, error)
	// 更新任务
	Update(ctx context.Context, task *Task) error
	// 删除任务
	Delete(ctx context.Context, id string) error
	// 获取任务列表
	List(ctx context.Context, workID, batchID string, taskType TaskType, status TaskStatus, page, size int) ([]*Task, int64, error)
	// 获取任务统计
	GetStats(ctx context.Context, workID string) (*TaskStats, error)
	// 更新任务状态
	UpdateStatus(ctx context.Context, id string, status TaskStatus) error
	// 增加重试次数
	IncrementRetryCount(ctx context.Context, id string) error
	// 获取待处理任务
	GetPendingTasks(ctx context.Context, limit int) ([]*Task, error)
	// 更新任务优先级
	UpdatePriority(ctx context.Context, id string, priority TaskPriority) error
	// 批量更新任务优先级
	BatchUpdatePriority(ctx context.Context, ids []string, priority TaskPriority) error
	// 根据条件更新任务优先级
	UpdatePriorityByCondition(ctx context.Context, workID, batchID string, taskType TaskType, status TaskStatus, priority TaskPriority) error

	// 优先级调整规则相关方法
	CreatePriorityRule(ctx context.Context, rule *PriorityAdjustRule) error
	GetPriorityRule(ctx context.Context, id string) (*PriorityAdjustRule, error)
	UpdatePriorityRule(ctx context.Context, rule *PriorityAdjustRule) error
	DeletePriorityRule(ctx context.Context, id string) error
	ListPriorityRules(ctx context.Context, workID, batchID string, taskType TaskType, enabled bool, page, size int) ([]*PriorityAdjustRule, int64, error)
	
	// 优先级调整日志相关方法
	CreatePriorityLog(ctx context.Context, log *PriorityAdjustLog) error
	GetPriorityLogs(ctx context.Context, taskID string, page, size int) ([]*PriorityAdjustLog, int64, error)
} 
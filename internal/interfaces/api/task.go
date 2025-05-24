package api

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"ai-translate/internal/infrastructure/ai"
	"ai-translate/internal/infrastructure/utils"
	"ai-translate/internal/model"
	"ai-translate/internal/service"
	"time"
)

// CreateTaskReq 创建任务请求
type CreateTaskReq struct {
	WorkID     string       `json:"work_id" v:"required#工作ID不能为空"`
	BatchID    string       `json:"batch_id" v:"required#批次ID不能为空"`
	Type       string       `json:"type" v:"required#任务类型不能为空"`
	Content    string       `json:"content" v:"required#内容不能为空"`
	Language   string       `json:"language"`
	SourceLang string       `json:"source_lang"`
	TargetLang string       `json:"target_lang"`
	Driver     ai.DriverType `json:"driver" v:"required#AI驱动不能为空"`
	Priority   int          `json:"priority" v:"min:0,max:3#优先级必须在0-3之间"`
}

// CreateTaskRes 创建任务响应
type CreateTaskRes struct {
	TaskID string `json:"task_id"`
}

// TaskController 任务控制器
type TaskController struct {
	taskService *service.TaskService
}

// NewTaskController 创建任务控制器
func NewTaskController() (*TaskController, error) {
	taskService, err := service.NewTaskService()
	if err != nil {
		return nil, err
	}
	return &TaskController{
		taskService: taskService,
	}, nil
}

// Create 创建任务
func (c *TaskController) Create(ctx context.Context, req *CreateTaskReq) (*CreateTaskRes, error) {
	// 创建任务
	task, err := c.taskService.CreateTask(ctx, req.WorkID, req.BatchID, model.TaskType(req.Type), req.Content, req.Driver, model.TaskPriority(req.Priority), req.Language, req.SourceLang, req.TargetLang)
	if err != nil {
		return nil, utils.NewError(utils.ErrInternalServer, "创建任务失败")
	}

	return &CreateTaskRes{
		TaskID: task.ID,
	}, nil
}

// GetTaskReq 获取任务请求
type GetTaskReq struct {
	TaskID string `json:"task_id" v:"required#任务ID不能为空"`
}

// GetTaskRes 获取任务响应
type GetTaskRes struct {
	TaskID     string       `json:"task_id"`
	WorkID     string       `json:"work_id"`
	BatchID    string       `json:"batch_id"`
	Type       string       `json:"type"`
	Status     string       `json:"status"`
	Priority   int          `json:"priority"`
	Content    string       `json:"content"`
	Result     string       `json:"result"`
	Error      string       `json:"error"`
	Driver     ai.DriverType `json:"driver"`
	RetryCount int          `json:"retry_count"`
	MaxRetries int          `json:"max_retries"`
	StartedAt  string       `json:"started_at"`
	CompletedAt string      `json:"completed_at"`
	CreatedAt  string       `json:"created_at"`
	UpdatedAt  string       `json:"updated_at"`
}

// Get 获取任务
func (c *TaskController) Get(ctx context.Context, req *GetTaskReq) (*GetTaskRes, error) {
	task, err := c.taskService.GetTask(ctx, req.TaskID)
	if err != nil {
		return nil, utils.NewError(utils.ErrNotFound, "任务不存在")
	}

	return &GetTaskRes{
		TaskID:     task.ID,
		WorkID:     task.WorkID,
		BatchID:    task.BatchID,
		Type:       string(task.Type),
		Status:     string(task.Status),
		Priority:   int(task.Priority),
		Content:    task.Content,
		Result:     task.Result,
		Error:      task.Error,
		Driver:     task.Driver,
		RetryCount: task.RetryCount,
		MaxRetries: task.MaxRetries,
		StartedAt:  task.StartedAt.Format("2006-01-02 15:04:05"),
		CompletedAt: task.CompletedAt.Format("2006-01-02 15:04:05"),
		CreatedAt:  task.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  task.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// ListTasksReq 获取任务列表请求
type ListTasksReq struct {
	WorkID  string `json:"work_id"`
	BatchID string `json:"batch_id"`
	Type    string `json:"type"`
	Status  string `json:"status"`
	Page    int    `json:"page" v:"min:1#页码必须大于0"`
	Size    int    `json:"size" v:"min:1,max:100#每页数量必须在1-100之间"`
}

// ListTasksRes 获取任务列表响应
type ListTasksRes struct {
	Total int64        `json:"total"`
	List  []*GetTaskRes `json:"list"`
}

// List 获取任务列表
func (c *TaskController) List(ctx context.Context, req *ListTasksReq) (*ListTasksRes, error) {
	tasks, total, err := c.taskService.ListTasks(ctx, req.WorkID, req.BatchID, model.TaskType(req.Type), model.TaskStatus(req.Status), req.Page, req.Size)
	if err != nil {
		return nil, utils.NewError(utils.ErrInternalServer, "获取任务列表失败")
	}

	list := make([]*GetTaskRes, 0, len(tasks))
	for _, task := range tasks {
		list = append(list, &GetTaskRes{
			TaskID:     task.ID,
			WorkID:     task.WorkID,
			BatchID:    task.BatchID,
			Type:       string(task.Type),
			Status:     string(task.Status),
			Priority:   int(task.Priority),
			Content:    task.Content,
			Result:     task.Result,
			Error:      task.Error,
			Driver:     task.Driver,
			RetryCount: task.RetryCount,
			MaxRetries: task.MaxRetries,
			StartedAt:  task.StartedAt.Format("2006-01-02 15:04:05"),
			CompletedAt: task.CompletedAt.Format("2006-01-02 15:04:05"),
			CreatedAt:  task.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  task.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &ListTasksRes{
		Total: total,
		List:  list,
	}, nil
}

// GetTaskStatsReq 获取任务统计请求
type GetTaskStatsReq struct {
	WorkID string `json:"work_id"`
}

// GetTaskStatsRes 获取任务统计响应
type GetTaskStatsRes struct {
	TotalTasks      int64   `json:"total_tasks"`
	PendingTasks    int64   `json:"pending_tasks"`
	RunningTasks    int64   `json:"running_tasks"`
	CompletedTasks  int64   `json:"completed_tasks"`
	FailedTasks     int64   `json:"failed_tasks"`
	PausedTasks     int64   `json:"paused_tasks"`
	AvgProcessTime  int64   `json:"avg_process_time"`
	SuccessRate     float64 `json:"success_rate"`
}

// GetStats 获取任务统计
func (c *TaskController) GetStats(ctx context.Context, req *GetTaskStatsReq) (*GetTaskStatsRes, error) {
	stats, err := c.taskService.GetTaskStats(ctx, req.WorkID)
	if err != nil {
		return nil, utils.NewError(utils.ErrInternalServer, "获取任务统计失败")
	}

	return &GetTaskStatsRes{
		TotalTasks:      stats.TotalTasks,
		PendingTasks:    stats.PendingTasks,
		RunningTasks:    stats.RunningTasks,
		CompletedTasks:  stats.CompletedTasks,
		FailedTasks:     stats.FailedTasks,
		PausedTasks:     stats.PausedTasks,
		AvgProcessTime:  stats.AvgProcessTime,
		SuccessRate:     stats.SuccessRate,
	}, nil
}

// UpdateTaskStatusReq 更新任务状态请求
type UpdateTaskStatusReq struct {
	TaskID string `json:"task_id" v:"required#任务ID不能为空"`
	Status string `json:"status" v:"required#任务状态不能为空"`
}

// UpdateTaskStatusRes 更新任务状态响应
type UpdateTaskStatusRes struct {
	Success bool `json:"success"`
}

// UpdateStatus 更新任务状态
func (c *TaskController) UpdateStatus(ctx context.Context, req *UpdateTaskStatusReq) (*UpdateTaskStatusRes, error) {
	if err := c.taskService.UpdateTaskStatus(ctx, req.TaskID, model.TaskStatus(req.Status)); err != nil {
		return nil, utils.NewError(utils.ErrInternalServer, "更新任务状态失败")
	}

	return &UpdateTaskStatusRes{
		Success: true,
	}, nil
}

// DeleteTaskReq 删除任务请求
type DeleteTaskReq struct {
	TaskID string `json:"task_id" v:"required#任务ID不能为空"`
}

// DeleteTaskRes 删除任务响应
type DeleteTaskRes struct {
	Success bool `json:"success"`
}

// Delete 删除任务
func (c *TaskController) Delete(ctx context.Context, req *DeleteTaskReq) (*DeleteTaskRes, error) {
	if err := c.taskService.DeleteTask(ctx, req.TaskID); err != nil {
		return nil, utils.NewError(utils.ErrInternalServer, "删除任务失败")
	}

	return &DeleteTaskRes{
		Success: true,
	}, nil
}

// RetryTaskReq 重试任务请求
type RetryTaskReq struct {
	TaskID string `json:"task_id" v:"required#任务ID不能为空"`
}

// RetryTaskRes 重试任务响应
type RetryTaskRes struct {
	Success bool `json:"success"`
}

// Retry 重试任务
func (c *TaskController) Retry(ctx context.Context, req *RetryTaskReq) (*RetryTaskRes, error) {
	if err := c.taskService.RetryTask(ctx, req.TaskID); err != nil {
		return nil, utils.NewError(utils.ErrInternalServer, "重试任务失败")
	}

	return &RetryTaskRes{
		Success: true,
	}, nil
}

// UpdateTaskPriorityReq 更新任务优先级请求
type UpdateTaskPriorityReq struct {
	TaskID   string `json:"task_id" v:"required#任务ID不能为空"`
	Priority int    `json:"priority" v:"min:0,max:3#优先级必须在0-3之间"`
}

// UpdateTaskPriorityRes 更新任务优先级响应
type UpdateTaskPriorityRes struct {
	Success bool `json:"success"`
}

// UpdatePriority 更新任务优先级
func (c *TaskController) UpdatePriority(ctx context.Context, req *UpdateTaskPriorityReq) (*UpdateTaskPriorityRes, error) {
	if err := c.taskService.UpdateTaskPriority(ctx, req.TaskID, model.TaskPriority(req.Priority)); err != nil {
		return nil, utils.NewError(utils.ErrInternalServer, "更新任务优先级失败")
	}

	return &UpdateTaskPriorityRes{
		Success: true,
	}, nil
}

// BatchUpdateTaskPriorityReq 批量更新任务优先级请求
type BatchUpdateTaskPriorityReq struct {
	TaskIDs  []string `json:"task_ids" v:"required#任务ID列表不能为空"`
	Priority int      `json:"priority" v:"min:0,max:3#优先级必须在0-3之间"`
}

// BatchUpdateTaskPriorityRes 批量更新任务优先级响应
type BatchUpdateTaskPriorityRes struct {
	Success bool `json:"success"`
}

// BatchUpdatePriority 批量更新任务优先级
func (c *TaskController) BatchUpdatePriority(ctx context.Context, req *BatchUpdateTaskPriorityReq) (*BatchUpdateTaskPriorityRes, error) {
	if err := c.taskService.BatchUpdateTaskPriority(ctx, req.TaskIDs, model.TaskPriority(req.Priority)); err != nil {
		return nil, utils.NewError(utils.ErrInternalServer, "批量更新任务优先级失败")
	}

	return &BatchUpdateTaskPriorityRes{
		Success: true,
	}, nil
}

// UpdatePriorityByConditionReq 根据条件更新任务优先级请求
type UpdatePriorityByConditionReq struct {
	WorkID  string `json:"work_id"`
	BatchID string `json:"batch_id"`
	Type    string `json:"type"`
	Status  string `json:"status"`
	Priority int    `json:"priority" v:"min:0,max:3#优先级必须在0-3之间"`
}

// UpdatePriorityByConditionRes 根据条件更新任务优先级响应
type UpdatePriorityByConditionRes struct {
	Success bool `json:"success"`
}

// UpdatePriorityByCondition 根据条件更新任务优先级
func (c *TaskController) UpdatePriorityByCondition(ctx context.Context, req *UpdatePriorityByConditionReq) (*UpdatePriorityByConditionRes, error) {
	if err := c.taskService.UpdateTaskPriorityByCondition(ctx, req.WorkID, req.BatchID, model.TaskType(req.Type), model.TaskStatus(req.Status), model.TaskPriority(req.Priority)); err != nil {
		return nil, utils.NewError(utils.ErrInternalServer, "更新任务优先级失败")
	}

	return &UpdatePriorityByConditionRes{
		Success: true,
	}, nil
}

// CreatePriorityRuleReq 创建优先级规则请求
type CreatePriorityRuleReq struct {
	WorkID      string         `json:"work_id" v:"required"`      // 工作ID
	BatchID     string         `json:"batch_id"`                  // 批次ID
	TaskType    model.TaskType `json:"task_type"`                 // 任务类型
	Status      model.TaskStatus `json:"status"`                  // 任务状态
	Priority    int            `json:"priority" v:"required"`     // 优先级
	Description string         `json:"description" v:"required"`  // 规则描述
	Enabled     bool           `json:"enabled"`                   // 是否启用
}

// CreatePriorityRuleRes 创建优先级规则响应
type CreatePriorityRuleRes struct {
	ID string `json:"id"` // 规则ID
}

// GetPriorityRuleReq 获取优先级规则请求
type GetPriorityRuleReq struct {
	ID string `json:"id" v:"required"` // 规则ID
}

// GetPriorityRuleRes 获取优先级规则响应
type GetPriorityRuleRes struct {
	Rule *model.PriorityAdjustRule `json:"rule"` // 规则信息
}

// UpdatePriorityRuleReq 更新优先级规则请求
type UpdatePriorityRuleReq struct {
	ID          string         `json:"id" v:"required"`          // 规则ID
	WorkID      string         `json:"work_id"`                  // 工作ID
	BatchID     string         `json:"batch_id"`                 // 批次ID
	TaskType    model.TaskType `json:"task_type"`                // 任务类型
	Status      model.TaskStatus `json:"status"`                 // 任务状态
	Priority    int            `json:"priority"`                 // 优先级
	Description string         `json:"description"`              // 规则描述
	Enabled     bool           `json:"enabled"`                  // 是否启用
}

// UpdatePriorityRuleRes 更新优先级规则响应
type UpdatePriorityRuleRes struct {
	Success bool `json:"success"` // 是否成功
}

// DeletePriorityRuleReq 删除优先级规则请求
type DeletePriorityRuleReq struct {
	ID string `json:"id" v:"required"` // 规则ID
}

// DeletePriorityRuleRes 删除优先级规则响应
type DeletePriorityRuleRes struct {
	Success bool `json:"success"` // 是否成功
}

// ListPriorityRulesReq 获取优先级规则列表请求
type ListPriorityRulesReq struct {
	WorkID   string         `json:"work_id"`    // 工作ID
	BatchID  string         `json:"batch_id"`   // 批次ID
	TaskType model.TaskType `json:"task_type"`  // 任务类型
	Enabled  bool           `json:"enabled"`    // 是否启用
	Page     int            `json:"page"`       // 页码
	Size     int            `json:"size"`       // 每页大小
}

// ListPriorityRulesRes 获取优先级规则列表响应
type ListPriorityRulesRes struct {
	Rules []*model.PriorityAdjustRule `json:"rules"` // 规则列表
	Total int64                       `json:"total"` // 总数
}

// GetPriorityLogsReq 获取优先级调整日志请求
type GetPriorityLogsReq struct {
	TaskID string `json:"task_id" v:"required"` // 任务ID
	Page   int    `json:"page"`                 // 页码
	Size   int    `json:"size"`                 // 每页大小
}

// GetPriorityLogsRes 获取优先级调整日志响应
type GetPriorityLogsRes struct {
	Logs  []*model.PriorityAdjustLog `json:"logs"`  // 日志列表
	Total int64                      `json:"total"` // 总数
}

// CreatePriorityRule 创建优先级规则
func (c *TaskController) CreatePriorityRule(ctx context.Context, req *CreatePriorityRuleReq) (*CreatePriorityRuleRes, error) {
	rule := &model.PriorityAdjustRule{
		ID:          utils.GenerateUUID(),
		WorkID:      req.WorkID,
		BatchID:     req.BatchID,
		TaskType:    req.TaskType,
		Status:      req.Status,
		Priority:    req.Priority,
		Description: req.Description,
		Enabled:     req.Enabled,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := c.taskService.CreatePriorityRule(ctx, rule); err != nil {
		return nil, err
	}

	return &CreatePriorityRuleRes{
		ID: rule.ID,
	}, nil
}

// GetPriorityRule 获取优先级规则
func (c *TaskController) GetPriorityRule(ctx context.Context, req *GetPriorityRuleReq) (*GetPriorityRuleRes, error) {
	rule, err := c.taskService.GetPriorityRule(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &GetPriorityRuleRes{
		Rule: rule,
	}, nil
}

// UpdatePriorityRule 更新优先级规则
func (c *TaskController) UpdatePriorityRule(ctx context.Context, req *UpdatePriorityRuleReq) (*UpdatePriorityRuleRes, error) {
	rule := &model.PriorityAdjustRule{
		ID:          req.ID,
		WorkID:      req.WorkID,
		BatchID:     req.BatchID,
		TaskType:    req.TaskType,
		Status:      req.Status,
		Priority:    req.Priority,
		Description: req.Description,
		Enabled:     req.Enabled,
		UpdatedAt:   time.Now(),
	}

	if err := c.taskService.UpdatePriorityRule(ctx, rule); err != nil {
		return nil, err
	}

	return &UpdatePriorityRuleRes{
		Success: true,
	}, nil
}

// DeletePriorityRule 删除优先级规则
func (c *TaskController) DeletePriorityRule(ctx context.Context, req *DeletePriorityRuleReq) (*DeletePriorityRuleRes, error) {
	if err := c.taskService.DeletePriorityRule(ctx, req.ID); err != nil {
		return nil, err
	}

	return &DeletePriorityRuleRes{
		Success: true,
	}, nil
}

// ListPriorityRules 获取优先级规则列表
func (c *TaskController) ListPriorityRules(ctx context.Context, req *ListPriorityRulesReq) (*ListPriorityRulesRes, error) {
	rules, total, err := c.taskService.ListPriorityRules(ctx, req.WorkID, req.BatchID, req.TaskType, req.Enabled, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	return &ListPriorityRulesRes{
		Rules: rules,
		Total: total,
	}, nil
}

// GetPriorityLogs 获取优先级调整日志
func (c *TaskController) GetPriorityLogs(ctx context.Context, req *GetPriorityLogsReq) (*GetPriorityLogsRes, error) {
	logs, total, err := c.taskService.GetPriorityLogs(ctx, req.TaskID, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	return &GetPriorityLogsRes{
		Logs:  logs,
		Total: total,
	}, nil
}

// CreateRuleTemplateReq 创建规则模板请求
type CreateRuleTemplateReq struct {
	Name        string            `json:"name" v:"required"`        // 模板名称
	Description string            `json:"description" v:"required"` // 模板描述
	Conditions  []model.RuleCondition `json:"conditions"`           // 条件列表
	Actions     []model.RuleAction    `json:"actions"`              // 动作列表
}

// CreateRuleTemplateRes 创建规则模板响应
type CreateRuleTemplateRes struct {
	ID string `json:"id"` // 模板ID
}

// GetRuleTemplateReq 获取规则模板请求
type GetRuleTemplateReq struct {
	ID string `json:"id" v:"required"` // 模板ID
}

// GetRuleTemplateRes 获取规则模板响应
type GetRuleTemplateRes struct {
	Template *model.RuleTemplate `json:"template"` // 模板信息
}

// UpdateRuleTemplateReq 更新规则模板请求
type UpdateRuleTemplateReq struct {
	ID          string            `json:"id" v:"required"`          // 模板ID
	Name        string            `json:"name"`                     // 模板名称
	Description string            `json:"description"`              // 模板描述
	Conditions  []model.RuleCondition `json:"conditions"`           // 条件列表
	Actions     []model.RuleAction    `json:"actions"`              // 动作列表
}

// UpdateRuleTemplateRes 更新规则模板响应
type UpdateRuleTemplateRes struct {
	Success bool `json:"success"` // 是否成功
}

// DeleteRuleTemplateReq 删除规则模板请求
type DeleteRuleTemplateReq struct {
	ID string `json:"id" v:"required"` // 模板ID
}

// DeleteRuleTemplateRes 删除规则模板响应
type DeleteRuleTemplateRes struct {
	Success bool `json:"success"` // 是否成功
}

// ListRuleTemplatesReq 获取规则模板列表请求
type ListRuleTemplatesReq struct {
	Page int `json:"page"` // 页码
	Size int `json:"size"` // 每页大小
}

// ListRuleTemplatesRes 获取规则模板列表响应
type ListRuleTemplatesRes struct {
	Templates []*model.RuleTemplate `json:"templates"` // 模板列表
	Total     int64                 `json:"total"`     // 总数
}

// CreateRuleGroupReq 创建规则组请求
type CreateRuleGroupReq struct {
	Name        string   `json:"name" v:"required"`        // 组名称
	Description string   `json:"description" v:"required"` // 组描述
	Rules       []string `json:"rules"`                    // 规则ID列表
	Enabled     bool     `json:"enabled"`                  // 是否启用
}

// CreateRuleGroupRes 创建规则组响应
type CreateRuleGroupRes struct {
	ID string `json:"id"` // 组ID
}

// GetRuleGroupReq 获取规则组请求
type GetRuleGroupReq struct {
	ID string `json:"id" v:"required"` // 组ID
}

// GetRuleGroupRes 获取规则组响应
type GetRuleGroupRes struct {
	Group *model.RuleGroup `json:"group"` // 组信息
}

// UpdateRuleGroupReq 更新规则组请求
type UpdateRuleGroupReq struct {
	ID          string   `json:"id" v:"required"`          // 组ID
	Name        string   `json:"name"`                     // 组名称
	Description string   `json:"description"`              // 组描述
	Rules       []string `json:"rules"`                    // 规则ID列表
	Enabled     bool     `json:"enabled"`                  // 是否启用
}

// UpdateRuleGroupRes 更新规则组响应
type UpdateRuleGroupRes struct {
	Success bool `json:"success"` // 是否成功
}

// DeleteRuleGroupReq 删除规则组请求
type DeleteRuleGroupReq struct {
	ID string `json:"id" v:"required"` // 组ID
}

// DeleteRuleGroupRes 删除规则组响应
type DeleteRuleGroupRes struct {
	Success bool `json:"success"` // 是否成功
}

// ListRuleGroupsReq 获取规则组列表请求
type ListRuleGroupsReq struct {
	Enabled bool `json:"enabled"` // 是否启用
	Page    int  `json:"page"`    // 页码
	Size    int  `json:"size"`    // 每页大小
}

// ListRuleGroupsRes 获取规则组列表响应
type ListRuleGroupsRes struct {
	Groups []*model.RuleGroup `json:"groups"` // 组列表
	Total  int64              `json:"total"`  // 总数
}

// AddRuleToGroupReq 添加规则到组请求
type AddRuleToGroupReq struct {
	GroupID string `json:"group_id" v:"required"` // 组ID
	RuleID  string `json:"rule_id" v:"required"`  // 规则ID
}

// AddRuleToGroupRes 添加规则到组响应
type AddRuleToGroupRes struct {
	Success bool `json:"success"` // 是否成功
}

// RemoveRuleFromGroupReq 从组中移除规则请求
type RemoveRuleFromGroupReq struct {
	GroupID string `json:"group_id" v:"required"` // 组ID
	RuleID  string `json:"rule_id" v:"required"`  // 规则ID
}

// RemoveRuleFromGroupRes 从组中移除规则响应
type RemoveRuleFromGroupRes struct {
	Success bool `json:"success"` // 是否成功
}

// GetGroupRulesReq 获取组中的规则请求
type GetGroupRulesReq struct {
	GroupID string `json:"group_id" v:"required"` // 组ID
}

// GetGroupRulesRes 获取组中的规则响应
type GetGroupRulesRes struct {
	Rules []*model.PriorityAdjustRule `json:"rules"` // 规则列表
}

// CreateRuleTemplate 创建规则模板
func (c *TaskController) CreateRuleTemplate(ctx context.Context, req *CreateRuleTemplateReq) (*CreateRuleTemplateRes, error) {
	template := &model.RuleTemplate{
		ID:          utils.GenerateUUID(),
		Name:        req.Name,
		Description: req.Description,
		Conditions:  req.Conditions,
		Actions:     req.Actions,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := c.taskService.CreateRuleTemplate(ctx, template); err != nil {
		return nil, err
	}

	return &CreateRuleTemplateRes{
		ID: template.ID,
	}, nil
}

// GetRuleTemplate 获取规则模板
func (c *TaskController) GetRuleTemplate(ctx context.Context, req *GetRuleTemplateReq) (*GetRuleTemplateRes, error) {
	template, err := c.taskService.GetRuleTemplate(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &GetRuleTemplateRes{
		Template: template,
	}, nil
}

// UpdateRuleTemplate 更新规则模板
func (c *TaskController) UpdateRuleTemplate(ctx context.Context, req *UpdateRuleTemplateReq) (*UpdateRuleTemplateRes, error) {
	template := &model.RuleTemplate{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Conditions:  req.Conditions,
		Actions:     req.Actions,
		UpdatedAt:   time.Now(),
	}

	if err := c.taskService.UpdateRuleTemplate(ctx, template); err != nil {
		return nil, err
	}

	return &UpdateRuleTemplateRes{
		Success: true,
	}, nil
}

// DeleteRuleTemplate 删除规则模板
func (c *TaskController) DeleteRuleTemplate(ctx context.Context, req *DeleteRuleTemplateReq) (*DeleteRuleTemplateRes, error) {
	if err := c.taskService.DeleteRuleTemplate(ctx, req.ID); err != nil {
		return nil, err
	}

	return &DeleteRuleTemplateRes{
		Success: true,
	}, nil
}

// ListRuleTemplates 获取规则模板列表
func (c *TaskController) ListRuleTemplates(ctx context.Context, req *ListRuleTemplatesReq) (*ListRuleTemplatesRes, error) {
	templates, total, err := c.taskService.ListRuleTemplates(ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	return &ListRuleTemplatesRes{
		Templates: templates,
		Total:     total,
	}, nil
}

// CreateRuleGroup 创建规则组
func (c *TaskController) CreateRuleGroup(ctx context.Context, req *CreateRuleGroupReq) (*CreateRuleGroupRes, error) {
	group := &model.RuleGroup{
		ID:          utils.GenerateUUID(),
		Name:        req.Name,
		Description: req.Description,
		Rules:       req.Rules,
		Enabled:     req.Enabled,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := c.taskService.CreateRuleGroup(ctx, group); err != nil {
		return nil, err
	}

	return &CreateRuleGroupRes{
		ID: group.ID,
	}, nil
}

// GetRuleGroup 获取规则组
func (c *TaskController) GetRuleGroup(ctx context.Context, req *GetRuleGroupReq) (*GetRuleGroupRes, error) {
	group, err := c.taskService.GetRuleGroup(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &GetRuleGroupRes{
		Group: group,
	}, nil
}

// UpdateRuleGroup 更新规则组
func (c *TaskController) UpdateRuleGroup(ctx context.Context, req *UpdateRuleGroupReq) (*UpdateRuleGroupRes, error) {
	group := &model.RuleGroup{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Rules:       req.Rules,
		Enabled:     req.Enabled,
		UpdatedAt:   time.Now(),
	}

	if err := c.taskService.UpdateRuleGroup(ctx, group); err != nil {
		return nil, err
	}

	return &UpdateRuleGroupRes{
		Success: true,
	}, nil
}

// DeleteRuleGroup 删除规则组
func (c *TaskController) DeleteRuleGroup(ctx context.Context, req *DeleteRuleGroupReq) (*DeleteRuleGroupRes, error) {
	if err := c.taskService.DeleteRuleGroup(ctx, req.ID); err != nil {
		return nil, err
	}

	return &DeleteRuleGroupRes{
		Success: true,
	}, nil
}

// ListRuleGroups 获取规则组列表
func (c *TaskController) ListRuleGroups(ctx context.Context, req *ListRuleGroupsReq) (*ListRuleGroupsRes, error) {
	groups, total, err := c.taskService.ListRuleGroups(ctx, req.Enabled, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	return &ListRuleGroupsRes{
		Groups: groups,
		Total:  total,
	}, nil
}

// AddRuleToGroup 添加规则到组
func (c *TaskController) AddRuleToGroup(ctx context.Context, req *AddRuleToGroupReq) (*AddRuleToGroupRes, error) {
	if err := c.taskService.AddRuleToGroup(ctx, req.GroupID, req.RuleID); err != nil {
		return nil, err
	}

	return &AddRuleToGroupRes{
		Success: true,
	}, nil
}

// RemoveRuleFromGroup 从组中移除规则
func (c *TaskController) RemoveRuleFromGroup(ctx context.Context, req *RemoveRuleFromGroupReq) (*RemoveRuleFromGroupRes, error) {
	if err := c.taskService.RemoveRuleFromGroup(ctx, req.GroupID, req.RuleID); err != nil {
		return nil, err
	}

	return &RemoveRuleFromGroupRes{
		Success: true,
	}, nil
}

// GetGroupRules 获取组中的规则
func (c *TaskController) GetGroupRules(ctx context.Context, req *GetGroupRulesReq) (*GetGroupRulesRes, error) {
	rules, err := c.taskService.GetGroupRules(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}

	return &GetGroupRulesRes{
		Rules: rules,
	}, nil
} 
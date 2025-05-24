package repository

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"gorm.io/gorm"
	"ai-translate/internal/model"
)

// TaskRepositoryImpl 任务仓储实现
type TaskRepositoryImpl struct {
	db *gorm.DB
}

// NewTaskRepository 创建任务仓储
func NewTaskRepository() (model.TaskRepository, error) {
	db, err := g.DB().GetDB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败: %v", err)
	}

	return &TaskRepositoryImpl{
		db: db,
	}, nil
}

// Create 创建任务
func (r *TaskRepositoryImpl) Create(ctx context.Context, task *model.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

// Get 获取任务
func (r *TaskRepositoryImpl) Get(ctx context.Context, id string) (*model.Task, error) {
	var task model.Task
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// Update 更新任务
func (r *TaskRepositoryImpl) Update(ctx context.Context, task *model.Task) error {
	return r.db.WithContext(ctx).Save(task).Error
}

// Delete 删除任务
func (r *TaskRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Task{}).Error
}

// List 获取任务列表
func (r *TaskRepositoryImpl) List(ctx context.Context, workID, batchID string, taskType model.TaskType, status model.TaskStatus, page, size int) ([]*model.Task, int64, error) {
	var tasks []*model.Task
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Task{})

	if workID != "" {
		query = query.Where("work_id = ?", workID)
	}
	if batchID != "" {
		query = query.Where("batch_id = ?", batchID)
	}
	if taskType != "" {
		query = query.Where("type = ?", taskType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := query.Offset((page - 1) * size).Limit(size).Order("priority DESC, created_at DESC").Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// GetStats 获取任务统计
func (r *TaskRepositoryImpl) GetStats(ctx context.Context, workID string) (*model.TaskStats, error) {
	var stats model.TaskStats
	query := r.db.WithContext(ctx).Model(&model.Task{})

	if workID != "" {
		query = query.Where("work_id = ?", workID)
	}

	// 获取各状态任务数量
	if err := query.Count(&stats.TotalTasks).Error; err != nil {
		return nil, err
	}
	if err := query.Where("status = ?", model.TaskStatusPending).Count(&stats.PendingTasks).Error; err != nil {
		return nil, err
	}
	if err := query.Where("status = ?", model.TaskStatusRunning).Count(&stats.RunningTasks).Error; err != nil {
		return nil, err
	}
	if err := query.Where("status = ?", model.TaskStatusCompleted).Count(&stats.CompletedTasks).Error; err != nil {
		return nil, err
	}
	if err := query.Where("status = ?", model.TaskStatusFailed).Count(&stats.FailedTasks).Error; err != nil {
		return nil, err
	}
	if err := query.Where("status = ?", model.TaskStatusPaused).Count(&stats.PausedTasks).Error; err != nil {
		return nil, err
	}

	// 计算平均处理时间
	var avgProcessTime int64
	if err := query.Where("status = ?", model.TaskStatusCompleted).
		Select("AVG(UNIX_TIMESTAMP(completed_at) - UNIX_TIMESTAMP(started_at))").
		Row().Scan(&avgProcessTime); err != nil {
		return nil, err
	}
	stats.AvgProcessTime = avgProcessTime

	// 计算成功率
	if stats.TotalTasks > 0 {
		stats.SuccessRate = float64(stats.CompletedTasks) / float64(stats.TotalTasks) * 100
	}

	return &stats, nil
}

// UpdateStatus 更新任务状态
func (r *TaskRepositoryImpl) UpdateStatus(ctx context.Context, id string, status model.TaskStatus) error {
	return r.db.WithContext(ctx).Model(&model.Task{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// IncrementRetryCount 增加重试次数
func (r *TaskRepositoryImpl) IncrementRetryCount(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.Task{}).
		Where("id = ?", id).
		UpdateColumn("retry_count", gorm.Expr("retry_count + ?", 1)).Error
}

// GetPendingTasks 获取待处理任务
func (r *TaskRepositoryImpl) GetPendingTasks(ctx context.Context, limit int) ([]*model.Task, error) {
	var tasks []*model.Task
	if err := r.db.WithContext(ctx).
		Where("status = ? AND retry_count < max_retries", model.TaskStatusPending).
		Order("priority DESC, created_at ASC").
		Limit(limit).
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// UpdatePriority 更新任务优先级
func (r *TaskRepositoryImpl) UpdatePriority(ctx context.Context, id string, priority model.TaskPriority) error {
	return r.db.WithContext(ctx).Model(&model.Task{}).
		Where("id = ?", id).
		Update("priority", priority).Error
}

// BatchUpdatePriority 批量更新任务优先级
func (r *TaskRepositoryImpl) BatchUpdatePriority(ctx context.Context, ids []string, priority model.TaskPriority) error {
	return r.db.WithContext(ctx).Model(&model.Task{}).
		Where("id IN ?", ids).
		Update("priority", priority).Error
}

// UpdatePriorityByCondition 根据条件更新任务优先级
func (r *TaskRepositoryImpl) UpdatePriorityByCondition(ctx context.Context, workID, batchID string, taskType model.TaskType, status model.TaskStatus, priority model.TaskPriority) error {
	query := r.db.WithContext(ctx).Model(&model.Task{})

	if workID != "" {
		query = query.Where("work_id = ?", workID)
	}
	if batchID != "" {
		query = query.Where("batch_id = ?", batchID)
	}
	if taskType != "" {
		query = query.Where("type = ?", taskType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	return query.Update("priority", priority).Error
}

// CreatePriorityRule 创建优先级调整规则
func (r *TaskRepositoryImpl) CreatePriorityRule(ctx context.Context, rule *model.PriorityAdjustRule) error {
	return r.db.WithContext(ctx).Create(rule).Error
}

// GetPriorityRule 获取优先级调整规则
func (r *TaskRepositoryImpl) GetPriorityRule(ctx context.Context, id string) (*model.PriorityAdjustRule, error) {
	var rule model.PriorityAdjustRule
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&rule).Error; err != nil {
		return nil, err
	}
	return &rule, nil
}

// UpdatePriorityRule 更新优先级调整规则
func (r *TaskRepositoryImpl) UpdatePriorityRule(ctx context.Context, rule *model.PriorityAdjustRule) error {
	return r.db.WithContext(ctx).Save(rule).Error
}

// DeletePriorityRule 删除优先级调整规则
func (r *TaskRepositoryImpl) DeletePriorityRule(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.PriorityAdjustRule{}).Error
}

// ListPriorityRules 获取优先级调整规则列表
func (r *TaskRepositoryImpl) ListPriorityRules(ctx context.Context, workID, batchID string, taskType model.TaskType, enabled bool, page, size int) ([]*model.PriorityAdjustRule, int64, error) {
	var rules []*model.PriorityAdjustRule
	var total int64

	query := r.db.WithContext(ctx).Model(&model.PriorityAdjustRule{})

	if workID != "" {
		query = query.Where("work_id = ?", workID)
	}
	if batchID != "" {
		query = query.Where("batch_id = ?", batchID)
	}
	if taskType != "" {
		query = query.Where("task_type = ?", taskType)
	}
	query = query.Where("enabled = ?", enabled)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := query.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&rules).Error; err != nil {
		return nil, 0, err
	}

	return rules, total, nil
}

// CreatePriorityLog 创建优先级调整日志
func (r *TaskRepositoryImpl) CreatePriorityLog(ctx context.Context, log *model.PriorityAdjustLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// GetPriorityLogs 获取优先级调整日志
func (r *TaskRepositoryImpl) GetPriorityLogs(ctx context.Context, taskID string, page, size int) ([]*model.PriorityAdjustLog, int64, error) {
	var logs []*model.PriorityAdjustLog
	var total int64

	query := r.db.WithContext(ctx).Model(&model.PriorityAdjustLog{}).Where("task_id = ?", taskID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := query.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// CreateRuleTemplate 创建规则模板
func (r *TaskRepositoryImpl) CreateRuleTemplate(ctx context.Context, template *model.RuleTemplate) error {
	return r.db.WithContext(ctx).Create(template).Error
}

// GetRuleTemplate 获取规则模板
func (r *TaskRepositoryImpl) GetRuleTemplate(ctx context.Context, id string) (*model.RuleTemplate, error) {
	var template model.RuleTemplate
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&template).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

// UpdateRuleTemplate 更新规则模板
func (r *TaskRepositoryImpl) UpdateRuleTemplate(ctx context.Context, template *model.RuleTemplate) error {
	return r.db.WithContext(ctx).Save(template).Error
}

// DeleteRuleTemplate 删除规则模板
func (r *TaskRepositoryImpl) DeleteRuleTemplate(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.RuleTemplate{}).Error
}

// ListRuleTemplates 获取规则模板列表
func (r *TaskRepositoryImpl) ListRuleTemplates(ctx context.Context, page, size int) ([]*model.RuleTemplate, int64, error) {
	var templates []*model.RuleTemplate
	var total int64

	query := r.db.WithContext(ctx).Model(&model.RuleTemplate{})

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := query.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&templates).Error; err != nil {
		return nil, 0, err
	}

	return templates, total, nil
}

// CreateRuleGroup 创建规则组
func (r *TaskRepositoryImpl) CreateRuleGroup(ctx context.Context, group *model.RuleGroup) error {
	return r.db.WithContext(ctx).Create(group).Error
}

// GetRuleGroup 获取规则组
func (r *TaskRepositoryImpl) GetRuleGroup(ctx context.Context, id string) (*model.RuleGroup, error) {
	var group model.RuleGroup
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&group).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

// UpdateRuleGroup 更新规则组
func (r *TaskRepositoryImpl) UpdateRuleGroup(ctx context.Context, group *model.RuleGroup) error {
	return r.db.WithContext(ctx).Save(group).Error
}

// DeleteRuleGroup 删除规则组
func (r *TaskRepositoryImpl) DeleteRuleGroup(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.RuleGroup{}).Error
}

// ListRuleGroups 获取规则组列表
func (r *TaskRepositoryImpl) ListRuleGroups(ctx context.Context, enabled bool, page, size int) ([]*model.RuleGroup, int64, error) {
	var groups []*model.RuleGroup
	var total int64

	query := r.db.WithContext(ctx).Model(&model.RuleGroup{}).Where("enabled = ?", enabled)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := query.Offset((page - 1) * size).Limit(size).Order("created_at DESC").Find(&groups).Error; err != nil {
		return nil, 0, err
	}

	return groups, total, nil
}

// AddRuleToGroup 添加规则到组
func (r *TaskRepositoryImpl) AddRuleToGroup(ctx context.Context, groupID, ruleID string) error {
	group, err := r.GetRuleGroup(ctx, groupID)
	if err != nil {
		return err
	}

	// 检查规则是否已存在
	for _, id := range group.Rules {
		if id == ruleID {
			return nil
		}
	}

	// 添加规则
	group.Rules = append(group.Rules, ruleID)
	return r.UpdateRuleGroup(ctx, group)
}

// RemoveRuleFromGroup 从组中移除规则
func (r *TaskRepositoryImpl) RemoveRuleFromGroup(ctx context.Context, groupID, ruleID string) error {
	group, err := r.GetRuleGroup(ctx, groupID)
	if err != nil {
		return err
	}

	// 移除规则
	newRules := make([]string, 0, len(group.Rules))
	for _, id := range group.Rules {
		if id != ruleID {
			newRules = append(newRules, id)
		}
	}
	group.Rules = newRules
	return r.UpdateRuleGroup(ctx, group)
}

// GetGroupRules 获取组中的规则
func (r *TaskRepositoryImpl) GetGroupRules(ctx context.Context, groupID string) ([]*model.PriorityAdjustRule, error) {
	group, err := r.GetRuleGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}

	var rules []*model.PriorityAdjustRule
	if err := r.db.WithContext(ctx).Where("id IN ?", group.Rules).Find(&rules).Error; err != nil {
		return nil, err
	}

	return rules, nil
}

// EvaluateRule 评估规则
func (r *TaskRepositoryImpl) EvaluateRule(ctx context.Context, rule *model.PriorityAdjustRule, task *model.Task) (bool, error) {
	// 检查工作ID
	if rule.WorkID != "" && rule.WorkID != task.WorkID {
		return false, nil
	}

	// 检查批次ID
	if rule.BatchID != "" && rule.BatchID != task.BatchID {
		return false, nil
	}

	// 检查任务类型
	if rule.TaskType != "" && rule.TaskType != task.Type {
		return false, nil
	}

	// 检查任务状态
	if rule.Status != "" && rule.Status != task.Status {
		return false, nil
	}

	return true, nil
}

// EvaluateRuleGroup 评估规则组
func (r *TaskRepositoryImpl) EvaluateRuleGroup(ctx context.Context, groupID string, task *model.Task) ([]*model.PriorityAdjustRule, error) {
	rules, err := r.GetGroupRules(ctx, groupID)
	if err != nil {
		return nil, err
	}

	var matchedRules []*model.PriorityAdjustRule
	for _, rule := range rules {
		matched, err := r.EvaluateRule(ctx, rule, task)
		if err != nil {
			return nil, err
		}
		if matched {
			matchedRules = append(matchedRules, rule)
		}
	}

	return matchedRules, nil
}

// GetApplicableRules 获取适用的规则
func (r *TaskRepositoryImpl) GetApplicableRules(ctx context.Context, task *model.Task) ([]*model.PriorityAdjustRule, error) {
	var rules []*model.PriorityAdjustRule
	if err := r.db.WithContext(ctx).
		Where("enabled = ?", true).
		Where("(work_id = ? OR work_id = '')", task.WorkID).
		Where("(batch_id = ? OR batch_id = '')", task.BatchID).
		Where("(task_type = ? OR task_type = '')", task.Type).
		Where("(status = ? OR status = '')", task.Status).
		Order("priority DESC").
		Find(&rules).Error; err != nil {
		return nil, err
	}

	return rules, nil
} 
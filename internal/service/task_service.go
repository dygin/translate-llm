package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"ai-translate/internal/infrastructure/ai"
	"ai-translate/internal/infrastructure/processor"
	"ai-translate/internal/infrastructure/queue"
	"ai-translate/internal/model"
	"time"
	"ai-translate/internal/utils"
)

// TaskService 任务服务
type TaskService struct {
	processor *processor.TaskProcessor
	aiDrivers map[ai.DriverType]ai.AIService
	repository *model.TaskRepository
}

// NewTaskService 创建任务服务
func NewTaskService() (*TaskService, error) {
	// 创建RabbitMQ客户端
	rabbitmq, err := queue.NewRabbitMQ()
	if err != nil {
		return nil, fmt.Errorf("创建RabbitMQ客户端失败: %v", err)
	}

	// 创建AI驱动
	aiDrivers := make(map[ai.DriverType]ai.AIService)
	
	// 创建OpenAI驱动
	openAIService, err := ai.NewAIService(ai.DriverOpenAI)
	if err != nil {
		return nil, fmt.Errorf("创建OpenAI服务失败: %v", err)
	}
	aiDrivers[ai.DriverOpenAI] = openAIService

	// 创建Gemini驱动
	geminiService, err := ai.NewAIService(ai.DriverGemini)
	if err != nil {
		return nil, fmt.Errorf("创建Gemini服务失败: %v", err)
	}
	aiDrivers[ai.DriverGemini] = geminiService

	// 获取配置
	workers := g.Cfg().MustGet("queue.worker.numWorkers").Int()
	maxRetries := g.Cfg().MustGet("queue.worker.maxRetries").Int()

	// 创建任务处理器
	taskProcessor := processor.NewTaskProcessor(rabbitmq, workers, maxRetries)

	// 创建任务服务
	service := &TaskService{
		processor: taskProcessor,
		aiDrivers: aiDrivers,
		repository: model.NewTaskRepository(),
	}

	// 注册任务处理函数
	taskProcessor.RegisterHandler("content_generation", service.handleContentGeneration)
	taskProcessor.RegisterHandler("translation", service.handleTranslation)

	return service, nil
}

// Start 启动任务服务
func (s *TaskService) Start(ctx context.Context) error {
	return s.processor.Start(ctx)
}

// Stop 停止任务服务
func (s *TaskService) Stop() {
	s.processor.Stop()
}

// CreatePriorityRule 创建优先级调整规则
func (s *TaskService) CreatePriorityRule(ctx context.Context, rule *model.PriorityAdjustRule) error {
	return s.repository.CreatePriorityRule(ctx, rule)
}

// GetPriorityRule 获取优先级调整规则
func (s *TaskService) GetPriorityRule(ctx context.Context, id string) (*model.PriorityAdjustRule, error) {
	return s.repository.GetPriorityRule(ctx, id)
}

// UpdatePriorityRule 更新优先级调整规则
func (s *TaskService) UpdatePriorityRule(ctx context.Context, rule *model.PriorityAdjustRule) error {
	return s.repository.UpdatePriorityRule(ctx, rule)
}

// DeletePriorityRule 删除优先级调整规则
func (s *TaskService) DeletePriorityRule(ctx context.Context, id string) error {
	return s.repository.DeletePriorityRule(ctx, id)
}

// ListPriorityRules 获取优先级调整规则列表
func (s *TaskService) ListPriorityRules(ctx context.Context, workID, batchID string, taskType model.TaskType, enabled bool, page, size int) ([]*model.PriorityAdjustRule, int64, error) {
	return s.repository.ListPriorityRules(ctx, workID, batchID, taskType, enabled, page, size)
}

// GetPriorityLogs 获取优先级调整日志
func (s *TaskService) GetPriorityLogs(ctx context.Context, taskID string, page, size int) ([]*model.PriorityAdjustLog, int64, error) {
	return s.repository.GetPriorityLogs(ctx, taskID, page, size)
}

// CreateRuleTemplate 创建规则模板
func (s *TaskService) CreateRuleTemplate(ctx context.Context, template *model.RuleTemplate) error {
	return s.repository.CreateRuleTemplate(ctx, template)
}

// GetRuleTemplate 获取规则模板
func (s *TaskService) GetRuleTemplate(ctx context.Context, id string) (*model.RuleTemplate, error) {
	return s.repository.GetRuleTemplate(ctx, id)
}

// UpdateRuleTemplate 更新规则模板
func (s *TaskService) UpdateRuleTemplate(ctx context.Context, template *model.RuleTemplate) error {
	return s.repository.UpdateRuleTemplate(ctx, template)
}

// DeleteRuleTemplate 删除规则模板
func (s *TaskService) DeleteRuleTemplate(ctx context.Context, id string) error {
	return s.repository.DeleteRuleTemplate(ctx, id)
}

// ListRuleTemplates 获取规则模板列表
func (s *TaskService) ListRuleTemplates(ctx context.Context, page, size int) ([]*model.RuleTemplate, int64, error) {
	return s.repository.ListRuleTemplates(ctx, page, size)
}

// CreateRuleGroup 创建规则组
func (s *TaskService) CreateRuleGroup(ctx context.Context, group *model.RuleGroup) error {
	return s.repository.CreateRuleGroup(ctx, group)
}

// GetRuleGroup 获取规则组
func (s *TaskService) GetRuleGroup(ctx context.Context, id string) (*model.RuleGroup, error) {
	return s.repository.GetRuleGroup(ctx, id)
}

// UpdateRuleGroup 更新规则组
func (s *TaskService) UpdateRuleGroup(ctx context.Context, group *model.RuleGroup) error {
	return s.repository.UpdateRuleGroup(ctx, group)
}

// DeleteRuleGroup 删除规则组
func (s *TaskService) DeleteRuleGroup(ctx context.Context, id string) error {
	return s.repository.DeleteRuleGroup(ctx, id)
}

// ListRuleGroups 获取规则组列表
func (s *TaskService) ListRuleGroups(ctx context.Context, enabled bool, page, size int) ([]*model.RuleGroup, int64, error) {
	return s.repository.ListRuleGroups(ctx, enabled, page, size)
}

// AddRuleToGroup 添加规则到组
func (s *TaskService) AddRuleToGroup(ctx context.Context, groupID, ruleID string) error {
	return s.repository.AddRuleToGroup(ctx, groupID, ruleID)
}

// RemoveRuleFromGroup 从组中移除规则
func (s *TaskService) RemoveRuleFromGroup(ctx context.Context, groupID, ruleID string) error {
	return s.repository.RemoveRuleFromGroup(ctx, groupID, ruleID)
}

// GetGroupRules 获取组中的规则
func (s *TaskService) GetGroupRules(ctx context.Context, groupID string) ([]*model.PriorityAdjustRule, error) {
	return s.repository.GetGroupRules(ctx, groupID)
}

// EvaluateRule 评估规则
func (s *TaskService) EvaluateRule(ctx context.Context, rule *model.PriorityAdjustRule, task *model.Task) (bool, error) {
	return s.repository.EvaluateRule(ctx, rule, task)
}

// EvaluateRuleGroup 评估规则组
func (s *TaskService) EvaluateRuleGroup(ctx context.Context, groupID string, task *model.Task) ([]*model.PriorityAdjustRule, error) {
	return s.repository.EvaluateRuleGroup(ctx, groupID, task)
}

// GetApplicableRules 获取适用的规则
func (s *TaskService) GetApplicableRules(ctx context.Context, task *model.Task) ([]*model.PriorityAdjustRule, error) {
	return s.repository.GetApplicableRules(ctx, task)
}

// adjustTaskPriority 调整任务优先级
func (s *TaskService) adjustTaskPriority(ctx context.Context, task *model.Task) error {
	// 获取适用的规则
	rules, err := s.GetApplicableRules(ctx, task)
	if err != nil {
		return fmt.Errorf("获取适用规则失败: %v", err)
	}

	// 遍历规则，找到第一个匹配的规则
	for _, rule := range rules {
		// 评估规则
		matched, err := s.EvaluateRule(ctx, rule, task)
		if err != nil {
			return fmt.Errorf("评估规则失败: %v", err)
		}

		if !matched {
			continue
		}

		// 如果优先级相同，不需要调整
		if rule.Priority == task.Priority {
			continue
		}

		// 记录旧优先级
		oldPriority := task.Priority

		// 更新任务优先级
		if err := s.repository.UpdatePriority(ctx, task.ID, rule.Priority); err != nil {
			return fmt.Errorf("更新任务优先级失败: %v", err)
		}

		// 创建优先级调整日志
		log := &model.PriorityAdjustLog{
			ID:          utils.GenerateUUID(),
			TaskID:      task.ID,
			RuleID:      rule.ID,
			OldPriority: oldPriority,
			NewPriority: rule.Priority,
			Reason:      rule.Description,
			CreatedAt:   time.Now(),
		}
		if err := s.repository.CreatePriorityLog(ctx, log); err != nil {
			return fmt.Errorf("创建优先级调整日志失败: %v", err)
		}

		// 记录日志
		g.Log().Infof(ctx, "任务优先级已调整: task_id=%s, old_priority=%d, new_priority=%d, reason=%s",
			task.ID, oldPriority, rule.Priority, rule.Description)

		return nil
	}

	return nil
}

// handleContentGeneration 处理内容生成任务
func (s *TaskService) handleContentGeneration(ctx context.Context, taskID string, data []byte) error {
	// 解析任务数据
	var taskData struct {
		WorkID    string `json:"work_id"`
		BatchID   string `json:"batch_id"`
		Content   string `json:"content"`
		Language  string `json:"language"`
		Driver    string `json:"driver"`
	}
	if err := json.Unmarshal(data, &taskData); err != nil {
		return fmt.Errorf("解析任务数据失败: %v", err)
	}

	// 获取AI驱动
	driver := ai.DriverType(taskData.Driver)
	aiService, ok := s.aiDrivers[driver]
	if !ok {
		return fmt.Errorf("不支持的AI驱动类型: %s", driver)
	}

	// 生成内容
	generatedContent, err := aiService.GenerateContent(ctx, taskData.Content, taskData.Language)
	if err != nil {
		return fmt.Errorf("生成内容失败: %v", err)
	}

	// 获取任务信息
	task, err := s.repository.Get(ctx, taskID)
	if err != nil {
		return fmt.Errorf("获取任务信息失败: %v", err)
	}

	// 调整任务优先级
	if err := s.adjustTaskPriority(ctx, task); err != nil {
		g.Log().Warningf(ctx, "调整任务优先级失败: %v", err)
	}

	// TODO: 保存生成的内容到数据库
	g.Log().Infof(ctx, "内容生成成功: %s", taskID)

	return nil
}

// handleTranslation 处理翻译任务
func (s *TaskService) handleTranslation(ctx context.Context, taskID string, data []byte) error {
	// 解析任务数据
	var taskData struct {
		WorkID     string `json:"work_id"`
		BatchID    string `json:"batch_id"`
		Content    string `json:"content"`
		SourceLang string `json:"source_lang"`
		TargetLang string `json:"target_lang"`
		Driver     string `json:"driver"`
	}
	if err := json.Unmarshal(data, &taskData); err != nil {
		return fmt.Errorf("解析任务数据失败: %v", err)
	}

	// 获取AI驱动
	driver := ai.DriverType(taskData.Driver)
	aiService, ok := s.aiDrivers[driver]
	if !ok {
		return fmt.Errorf("不支持的AI驱动类型: %s", driver)
	}

	// 翻译内容
	translatedContent, err := aiService.Translate(ctx, taskData.Content, taskData.SourceLang, taskData.TargetLang)
	if err != nil {
		return fmt.Errorf("翻译内容失败: %v", err)
	}

	// 获取任务信息
	task, err := s.repository.Get(ctx, taskID)
	if err != nil {
		return fmt.Errorf("获取任务信息失败: %v", err)
	}

	// 调整任务优先级
	if err := s.adjustTaskPriority(ctx, task); err != nil {
		g.Log().Warningf(ctx, "调整任务优先级失败: %v", err)
	}

	// TODO: 保存翻译结果到数据库
	g.Log().Infof(ctx, "翻译成功: %s", taskID)

	return nil
}

// UpdateTaskPriority 更新任务优先级
func (s *TaskService) UpdateTaskPriority(ctx context.Context, taskID string, priority model.TaskPriority) error {
	return s.repository.UpdatePriority(ctx, taskID, priority)
}

// BatchUpdateTaskPriority 批量更新任务优先级
func (s *TaskService) BatchUpdateTaskPriority(ctx context.Context, taskIDs []string, priority model.TaskPriority) error {
	return s.repository.BatchUpdatePriority(ctx, taskIDs, priority)
}

// UpdateTaskPriorityByCondition 根据条件更新任务优先级
func (s *TaskService) UpdateTaskPriorityByCondition(ctx context.Context, workID, batchID string, taskType model.TaskType, status model.TaskStatus, priority model.TaskPriority) error {
	return s.repository.UpdatePriorityByCondition(ctx, workID, batchID, taskType, status, priority)
} 
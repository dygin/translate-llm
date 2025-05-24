package task

import (
	"ai-translate/internal/application"
	"ai-translate/internal/domain/ai"
	"ai-translate/internal/domain/task"
	"ai-translate/internal/domain/work"
	"ai-translate/internal/infrastructure/ai"
	"ai-translate/internal/infrastructure/storage"
	"context"
	"encoding/json"
	"github.com/gogf/gf/v2/frame/g"
	"time"
)

type Processor struct {
	taskService    task.TaskService
	workService    work.WorkService
	aiService      ai.ModelService
	storageService *storage.OSSService
}

// NewProcessor 创建任务处理器实例
func NewProcessor() (*Processor, error) {
	storageService, err := storage.NewOSSService()
	if err != nil {
		return nil, err
	}

	aiConfig := &ai.ModelConfig{
		Type:      ai.ModelTypeGemini,
		APIKey:    g.Cfg().MustGet(context.Background(), "gemini.apiKey").String(),
		ModelName: g.Cfg().MustGet(context.Background(), "gemini.model").String(),
	}

	aiService, err := ai.NewGeminiService(aiConfig)
	if err != nil {
		return nil, err
	}

	workService, err := application.NewWorkService()
	if err != nil {
		return nil, err
	}

	return &Processor{
		taskService:    application.NewTaskService(),
		workService:    workService,
		aiService:      aiService,
		storageService: storageService,
	}, nil
}

// Start 启动任务处理器
func (p *Processor) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			// 从队列中获取任务
			t, err := p.taskService.(*application.TaskService).TaskQueue.Pop()
			if err != nil {
				continue
			}

			// 处理任务
			err = p.processTask(ctx, t)
			if err != nil {
				// 更新任务状态为失败
				t.Status = 3 // 3:失败
				t.UpdatedAt = time.Now()
				p.taskService.UpdateTask(t)

				// 如果未超过最大重试次数，则重试
				if t.RetryCount < t.MaxRetry {
					p.taskService.RetryTask(t.ID)
				}
			} else {
				// 更新任务状态为完成
				t.Status = 1 // 1:完成
				t.UpdatedAt = time.Now()
				p.taskService.UpdateTask(t)
			}
		}
	}
}

// processTask 处理任务
func (p *Processor) processTask(ctx context.Context, t *task.Task) error {
	switch t.Type {
	case 1: // 内容生成
		return p.processContentGenerationTask(ctx, t)
	case 2: // 翻译
		return p.processTranslationTask(ctx, t)
	default:
		return nil
	}
}

// processContentGenerationTask 处理内容生成任务
func (p *Processor) processContentGenerationTask(ctx context.Context, t *task.Task) error {
	// 获取作品信息
	w, err := p.workService.GetWork(t.ReferenceID)
	if err != nil {
		return err
	}

	// 生成内容简介
	return p.workService.GenerateContentSummary(w.ID)
}

// processTranslationTask 处理翻译任务
func (p *Processor) processTranslationTask(ctx context.Context, t *task.Task) error {
	// 获取翻译批次信息
	batch, err := p.workService.GetTranslationBatch(t.ReferenceID)
	if err != nil {
		return err
	}

	// 获取作品信息
	w, err := p.workService.GetWork(batch.WorkID)
	if err != nil {
		return err
	}

	// 获取内容简介
	summary, err := p.workService.(*application.WorkService).ContentSummaryRepo.FindByWorkID(w.ID)
	if err != nil {
		return err
	}

	// 获取翻译提示词
	prompts, err := p.workService.(*application.WorkService).PromptRepo.FindByType(2) // 2:翻译
	if err != nil {
		return err
	}
	if len(prompts) == 0 {
		return nil
	}

	// 调用AI进行翻译
	req := &ai.TranslationRequest{
		Content:        summary.Content,
		TargetLanguage: batch.TargetLanguage,
		Terminology:    batch.TerminologyURL,
		Prompt:         prompts[0].Content,
	}

	resp, err := p.aiService.Translate(ctx, req)
	if err != nil {
		return err
	}

	// 生成SRT文件
	srtContent := generateSRT(resp.TranslatedContent)

	// 上传SRT文件到OSS
	objectKey := p.storageService.GenerateObjectKey("translations", w.Title)
	srtURL, err := p.storageService.UploadContent(objectKey, []byte(srtContent))
	if err != nil {
		return err
	}

	// 保存翻译结果
	result := &work.TranslationResult{
		BatchID:   batch.ID,
		SrtURL:    srtURL,
		Status:    1,
		CreatedAt: time.Now(),
	}

	return p.workService.(*application.WorkService).TranslationResultRepo.Save(result)
}

// generateSRT 生成SRT文件内容
func generateSRT(content string) string {
	// TODO: 实现SRT文件生成逻辑
	return content
} 
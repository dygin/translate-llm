package application

import (
	"ai-translate/internal/domain/ai"
	"ai-translate/internal/domain/prompt"
	"ai-translate/internal/domain/task"
	"ai-translate/internal/domain/work"
	"ai-translate/internal/infrastructure/ai"
	"ai-translate/internal/infrastructure/persistence"
	"ai-translate/internal/infrastructure/storage"
	"context"
	"errors"
	"github.com/gogf/gf/v2/frame/g"
	"time"
)

type workService struct {
	workRepo              work.WorkRepository
	contentSummaryRepo    work.ContentSummaryRepository
	translationBatchRepo  work.TranslationBatchRepository
	translationResultRepo work.TranslationResultRepository
	promptRepo           prompt.PromptRepository
	taskRepo             task.TaskRepository
	taskQueue            task.TaskQueue
	aiService            ai.ModelService
	storageService       *storage.OSSService
}

// NewWorkService 创建作品服务实例
func NewWorkService() (work.WorkService, error) {
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

	return &workService{
		workRepo:              persistence.NewWorkRepository(),
		contentSummaryRepo:    persistence.NewContentSummaryRepository(),
		translationBatchRepo:  persistence.NewTranslationBatchRepository(),
		translationResultRepo: persistence.NewTranslationResultRepository(),
		promptRepo:           persistence.NewPromptRepository(),
		taskRepo:             persistence.NewTaskRepository(),
		taskQueue:            persistence.NewTaskQueue(),
		aiService:            aiService,
		storageService:       storageService,
	}, nil
}

func (s *workService) CreateWork(work *work.Work) error {
	// 保存作品
	err := s.workRepo.Save(work)
	if err != nil {
		return err
	}

	// 创建内容生成任务
	task := &task.Task{
		Type:        1, // 内容生成
		Priority:    0,
		Status:      0,
		ReferenceID: work.ID,
		RetryCount:  0,
		MaxRetry:    3,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = s.taskRepo.Save(task)
	if err != nil {
		return err
	}

	// 将任务加入队列
	return s.taskQueue.Push(task)
}

func (s *workService) GetWork(id uint64) (*work.Work, error) {
	return s.workRepo.FindByID(id)
}

func (s *workService) GetUserWorks(userID uint64) ([]*work.Work, error) {
	return s.workRepo.FindByUserID(userID)
}

func (s *workService) UpdateWork(work *work.Work) error {
	return s.workRepo.Update(work)
}

func (s *workService) DeleteWork(id uint64) error {
	return s.workRepo.Delete(id)
}

func (s *workService) GenerateContentSummary(workID uint64) error {
	// 获取作品信息
	w, err := s.workRepo.FindByID(workID)
	if err != nil {
		return err
	}

	// 获取内容简介提示词
	prompts, err := s.promptRepo.FindByType(1) // 1:内容简介
	if err != nil {
		return err
	}
	if len(prompts) == 0 {
		return errors.New("未找到内容简介提示词")
	}

	// 调用AI生成内容
	req := &ai.GenerateContentRequest{
		VideoURL:    w.VideoURL,
		SubtitleURL: w.SubtitleURL,
		Prompt:      prompts[0].Content,
	}

	resp, err := s.aiService.GenerateContent(context.Background(), req)
	if err != nil {
		return err
	}

	// 上传内容到OSS
	objectKey := s.storageService.GenerateObjectKey("content_summaries", w.Title)
	ossURL, err := s.storageService.UploadContent(objectKey, []byte(resp.Content))
	if err != nil {
		return err
	}

	// 保存内容简介
	summary := &work.ContentSummary{
		WorkID:    workID,
		Content:   resp.Content,
		OssURL:    ossURL,
		Status:    1,
		CreatedAt: time.Now(),
	}

	return s.contentSummaryRepo.Save(summary)
}

func (s *workService) CreateTranslationBatch(batch *work.TranslationBatch) error {
	// 保存翻译批次
	err := s.translationBatchRepo.Save(batch)
	if err != nil {
		return err
	}

	// 创建翻译任务
	task := &task.Task{
		Type:        2, // 翻译
		Priority:    0,
		Status:      0,
		ReferenceID: batch.ID,
		RetryCount:  0,
		MaxRetry:    3,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = s.taskRepo.Save(task)
	if err != nil {
		return err
	}

	// 将任务加入队列
	return s.taskQueue.Push(task)
}

func (s *workService) GetTranslationBatch(id uint64) (*work.TranslationBatch, error) {
	return s.translationBatchRepo.FindByID(id)
}

func (s *workService) GetWorkTranslationBatches(workID uint64) ([]*work.TranslationBatch, error) {
	return s.translationBatchRepo.FindByWorkID(workID)
} 
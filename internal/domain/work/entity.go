package work

import (
	"time"
)

// Work 作品实体
type Work struct {
	ID          uint64    `json:"id"`
	Title       string    `json:"title"`
	UserID      uint64    `json:"user_id"`
	VideoURL    string    `json:"video_url"`
	SubtitleURL string    `json:"subtitle_url"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ContentSummary 内容简介实体
type ContentSummary struct {
	ID        uint64    `json:"id"`
	WorkID    uint64    `json:"work_id"`
	Content   string    `json:"content"`
	OssURL    string    `json:"oss_url"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// TranslationBatch 翻译批次实体
type TranslationBatch struct {
	ID             uint64    `json:"id"`
	WorkID         uint64    `json:"work_id"`
	TargetLanguage string    `json:"target_language"`
	TerminologyURL string    `json:"terminology_url"`
	Status         int       `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TranslationResult 翻译结果实体
type TranslationResult struct {
	ID        uint64    `json:"id"`
	BatchID   uint64    `json:"batch_id"`
	SrtURL    string    `json:"srt_url"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// WorkRepository 作品仓储接口
type WorkRepository interface {
	FindByID(id uint64) (*Work, error)
	FindByUserID(userID uint64) ([]*Work, error)
	Save(work *Work) error
	Update(work *Work) error
	Delete(id uint64) error
}

// ContentSummaryRepository 内容简介仓储接口
type ContentSummaryRepository interface {
	FindByWorkID(workID uint64) (*ContentSummary, error)
	Save(summary *ContentSummary) error
	Update(summary *ContentSummary) error
}

// TranslationBatchRepository 翻译批次仓储接口
type TranslationBatchRepository interface {
	FindByID(id uint64) (*TranslationBatch, error)
	FindByWorkID(workID uint64) ([]*TranslationBatch, error)
	Save(batch *TranslationBatch) error
	Update(batch *TranslationBatch) error
}

// TranslationResultRepository 翻译结果仓储接口
type TranslationResultRepository interface {
	FindByBatchID(batchID uint64) (*TranslationResult, error)
	Save(result *TranslationResult) error
	Update(result *TranslationResult) error
}

// WorkService 作品服务接口
type WorkService interface {
	CreateWork(work *Work) error
	GetWork(id uint64) (*Work, error)
	GetUserWorks(userID uint64) ([]*Work, error)
	UpdateWork(work *Work) error
	DeleteWork(id uint64) error
	GenerateContentSummary(workID uint64) error
	CreateTranslationBatch(batch *TranslationBatch) error
	GetTranslationBatch(id uint64) (*TranslationBatch, error)
	GetWorkTranslationBatches(workID uint64) ([]*TranslationBatch, error)
} 
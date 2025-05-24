package prompt

import (
	"time"
)

// Prompt 提示词实体
type Prompt struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	Type      int       `json:"type"` // 1:内容简介 2:翻译
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PromptRepository 提示词仓储接口
type PromptRepository interface {
	FindByID(id uint64) (*Prompt, error)
	FindByType(promptType int) ([]*Prompt, error)
	Save(prompt *Prompt) error
	Update(prompt *Prompt) error
	Delete(id uint64) error
}

// PromptService 提示词服务接口
type PromptService interface {
	CreatePrompt(prompt *Prompt) error
	GetPrompt(id uint64) (*Prompt, error)
	GetPromptsByType(promptType int) ([]*Prompt, error)
	UpdatePrompt(prompt *Prompt) error
	DeletePrompt(id uint64) error
} 
package application

import (
	"ai-translate/internal/domain/prompt"
	"ai-translate/internal/infrastructure/persistence"
)

type promptService struct {
	promptRepo prompt.PromptRepository
}

// NewPromptService 创建提示词服务实例
func NewPromptService() prompt.PromptService {
	return &promptService{
		promptRepo: persistence.NewPromptRepository(),
	}
}

func (s *promptService) CreatePrompt(prompt *prompt.Prompt) error {
	return s.promptRepo.Save(prompt)
}

func (s *promptService) GetPrompt(id uint64) (*prompt.Prompt, error) {
	return s.promptRepo.FindByID(id)
}

func (s *promptService) GetPromptsByType(promptType int) ([]*prompt.Prompt, error) {
	return s.promptRepo.FindByType(promptType)
}

func (s *promptService) UpdatePrompt(prompt *prompt.Prompt) error {
	return s.promptRepo.Update(prompt)
}

func (s *promptService) DeletePrompt(id uint64) error {
	return s.promptRepo.Delete(id)
} 
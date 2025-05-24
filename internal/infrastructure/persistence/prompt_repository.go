package persistence

import (
	"ai-translate/internal/domain/prompt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type promptRepository struct {
	db gdb.DB
}

// NewPromptRepository 创建提示词仓储实例
func NewPromptRepository() prompt.PromptRepository {
	return &promptRepository{
		db: g.DB(),
	}
}

func (r *promptRepository) FindByID(id uint64) (*prompt.Prompt, error) {
	var p prompt.Prompt
	err := r.db.Model("prompts").Where("id", id).Scan(&p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *promptRepository) FindByType(promptType int) ([]*prompt.Prompt, error) {
	var prompts []*prompt.Prompt
	err := r.db.Model("prompts").Where("type", promptType).Scan(&prompts)
	if err != nil {
		return nil, err
	}
	return prompts, nil
}

func (r *promptRepository) Save(prompt *prompt.Prompt) error {
	_, err := r.db.Model("prompts").Insert(prompt)
	return err
}

func (r *promptRepository) Update(prompt *prompt.Prompt) error {
	_, err := r.db.Model("prompts").Where("id", prompt.ID).Update(prompt)
	return err
}

func (r *promptRepository) Delete(id uint64) error {
	_, err := r.db.Model("prompts").Where("id", id).Delete()
	return err
} 
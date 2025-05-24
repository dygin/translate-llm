package persistence

import (
	"ai-translate/internal/domain/work"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type workRepository struct {
	db gdb.DB
}

// NewWorkRepository 创建作品仓储实例
func NewWorkRepository() work.WorkRepository {
	return &workRepository{
		db: g.DB(),
	}
}

func (r *workRepository) FindByID(id uint64) (*work.Work, error) {
	var w work.Work
	err := r.db.Model("works").Where("id", id).Scan(&w)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *workRepository) FindByUserID(userID uint64) ([]*work.Work, error) {
	var works []*work.Work
	err := r.db.Model("works").Where("user_id", userID).Scan(&works)
	if err != nil {
		return nil, err
	}
	return works, nil
}

func (r *workRepository) Save(work *work.Work) error {
	_, err := r.db.Model("works").Insert(work)
	return err
}

func (r *workRepository) Update(work *work.Work) error {
	_, err := r.db.Model("works").Where("id", work.ID).Update(work)
	return err
}

func (r *workRepository) Delete(id uint64) error {
	_, err := r.db.Model("works").Where("id", id).Delete()
	return err
}

type contentSummaryRepository struct {
	db gdb.DB
}

// NewContentSummaryRepository 创建内容简介仓储实例
func NewContentSummaryRepository() work.ContentSummaryRepository {
	return &contentSummaryRepository{
		db: g.DB(),
	}
}

func (r *contentSummaryRepository) FindByWorkID(workID uint64) (*work.ContentSummary, error) {
	var summary work.ContentSummary
	err := r.db.Model("content_summaries").Where("work_id", workID).Scan(&summary)
	if err != nil {
		return nil, err
	}
	return &summary, nil
}

func (r *contentSummaryRepository) Save(summary *work.ContentSummary) error {
	_, err := r.db.Model("content_summaries").Insert(summary)
	return err
}

func (r *contentSummaryRepository) Update(summary *work.ContentSummary) error {
	_, err := r.db.Model("content_summaries").Where("id", summary.ID).Update(summary)
	return err
}

type translationBatchRepository struct {
	db gdb.DB
}

// NewTranslationBatchRepository 创建翻译批次仓储实例
func NewTranslationBatchRepository() work.TranslationBatchRepository {
	return &translationBatchRepository{
		db: g.DB(),
	}
}

func (r *translationBatchRepository) FindByID(id uint64) (*work.TranslationBatch, error) {
	var batch work.TranslationBatch
	err := r.db.Model("translation_batches").Where("id", id).Scan(&batch)
	if err != nil {
		return nil, err
	}
	return &batch, nil
}

func (r *translationBatchRepository) FindByWorkID(workID uint64) ([]*work.TranslationBatch, error) {
	var batches []*work.TranslationBatch
	err := r.db.Model("translation_batches").Where("work_id", workID).Scan(&batches)
	if err != nil {
		return nil, err
	}
	return batches, nil
}

func (r *translationBatchRepository) Save(batch *work.TranslationBatch) error {
	_, err := r.db.Model("translation_batches").Insert(batch)
	return err
}

func (r *translationBatchRepository) Update(batch *work.TranslationBatch) error {
	_, err := r.db.Model("translation_batches").Where("id", batch.ID).Update(batch)
	return err
}

type translationResultRepository struct {
	db gdb.DB
}

// NewTranslationResultRepository 创建翻译结果仓储实例
func NewTranslationResultRepository() work.TranslationResultRepository {
	return &translationResultRepository{
		db: g.DB(),
	}
}

func (r *translationResultRepository) FindByBatchID(batchID uint64) (*work.TranslationResult, error) {
	var result work.TranslationResult
	err := r.db.Model("translation_results").Where("batch_id", batchID).Scan(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *translationResultRepository) Save(result *work.TranslationResult) error {
	_, err := r.db.Model("translation_results").Insert(result)
	return err
}

func (r *translationResultRepository) Update(result *work.TranslationResult) error {
	_, err := r.db.Model("translation_results").Where("id", result.ID).Update(result)
	return err
} 
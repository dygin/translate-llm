package api

import (
	"ai-translate/internal/application"
	"ai-translate/internal/domain/work"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type WorkController struct {
	workService work.WorkService
}

// NewWorkController 创建作品控制器实例
func NewWorkController() (*WorkController, error) {
	workService, err := application.NewWorkService()
	if err != nil {
		return nil, err
	}

	return &WorkController{
		workService: workService,
	}, nil
}

// CreateWork 创建作品
func (c *WorkController) CreateWork(r *ghttp.Request) {
	var req struct {
		Title       string `json:"title" v:"required"`
		VideoURL    string `json:"video_url" v:"required"`
		SubtitleURL string `json:"subtitle_url" v:"required"`
	}

	if err := r.Parse(&req); err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 400,
			"msg":  err.Error(),
		})
	}

	userID := r.GetCtxVar("user_id").Uint64()
	w := &work.Work{
		Title:       req.Title,
		UserID:      userID,
		VideoURL:    req.VideoURL,
		SubtitleURL: req.SubtitleURL,
		Status:      0,
	}

	err := c.workService.CreateWork(w)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}

	r.Response.WriteJsonExit(g.Map{
		"code": 200,
		"msg":  "创建成功",
		"data": w,
	})
}

// GetWork 获取作品信息
func (c *WorkController) GetWork(r *ghttp.Request) {
	id := r.Get("id").Uint64()
	w, err := c.workService.GetWork(id)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}

	r.Response.WriteJsonExit(g.Map{
		"code": 200,
		"msg":  "获取成功",
		"data": w,
	})
}

// GetUserWorks 获取用户作品列表
func (c *WorkController) GetUserWorks(r *ghttp.Request) {
	userID := r.GetCtxVar("user_id").Uint64()
	works, err := c.workService.GetUserWorks(userID)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}

	r.Response.WriteJsonExit(g.Map{
		"code": 200,
		"msg":  "获取成功",
		"data": works,
	})
}

// UpdateWork 更新作品信息
func (c *WorkController) UpdateWork(r *ghttp.Request) {
	var req struct {
		ID          uint64 `json:"id" v:"required"`
		Title       string `json:"title"`
		VideoURL    string `json:"video_url"`
		SubtitleURL string `json:"subtitle_url"`
	}

	if err := r.Parse(&req); err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 400,
			"msg":  err.Error(),
		})
	}

	w, err := c.workService.GetWork(req.ID)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}

	if req.Title != "" {
		w.Title = req.Title
	}
	if req.VideoURL != "" {
		w.VideoURL = req.VideoURL
	}
	if req.SubtitleURL != "" {
		w.SubtitleURL = req.SubtitleURL
	}

	err = c.workService.UpdateWork(w)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}

	r.Response.WriteJsonExit(g.Map{
		"code": 200,
		"msg":  "更新成功",
		"data": w,
	})
}

// DeleteWork 删除作品
func (c *WorkController) DeleteWork(r *ghttp.Request) {
	id := r.Get("id").Uint64()
	err := c.workService.DeleteWork(id)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}

	r.Response.WriteJsonExit(g.Map{
		"code": 200,
		"msg":  "删除成功",
	})
}

// CreateTranslationBatch 创建翻译批次
func (c *WorkController) CreateTranslationBatch(r *ghttp.Request) {
	var req struct {
		WorkID         uint64 `json:"work_id" v:"required"`
		TargetLanguage string `json:"target_language" v:"required"`
		TerminologyURL string `json:"terminology_url"`
	}

	if err := r.Parse(&req); err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 400,
			"msg":  err.Error(),
		})
	}

	batch := &work.TranslationBatch{
		WorkID:         req.WorkID,
		TargetLanguage: req.TargetLanguage,
		TerminologyURL: req.TerminologyURL,
		Status:         0,
	}

	err := c.workService.CreateTranslationBatch(batch)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}

	r.Response.WriteJsonExit(g.Map{
		"code": 200,
		"msg":  "创建成功",
		"data": batch,
	})
}

// GetTranslationBatch 获取翻译批次信息
func (c *WorkController) GetTranslationBatch(r *ghttp.Request) {
	id := r.Get("id").Uint64()
	batch, err := c.workService.GetTranslationBatch(id)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}

	r.Response.WriteJsonExit(g.Map{
		"code": 200,
		"msg":  "获取成功",
		"data": batch,
	})
}

// GetWorkTranslationBatches 获取作品翻译批次列表
func (c *WorkController) GetWorkTranslationBatches(r *ghttp.Request) {
	workID := r.Get("work_id").Uint64()
	batches, err := c.workService.GetWorkTranslationBatches(workID)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}

	r.Response.WriteJsonExit(g.Map{
		"code": 200,
		"msg":  "获取成功",
		"data": batches,
	})
} 
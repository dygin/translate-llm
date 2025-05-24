package api

import (
	"ai-translate/internal/application"
	"ai-translate/internal/domain/prompt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type PromptController struct {
	promptService prompt.PromptService
}

// NewPromptController 创建提示词控制器实例
func NewPromptController() *PromptController {
	return &PromptController{
		promptService: application.NewPromptService(),
	}
}

// CreatePrompt 创建提示词
func (c *PromptController) CreatePrompt(r *ghttp.Request) {
	var req struct {
		Name    string `json:"name" v:"required"`
		Type    int    `json:"type" v:"required|in:1,2"`
		Content string `json:"content" v:"required"`
	}

	if err := r.Parse(&req); err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 400,
			"msg":  err.Error(),
		})
	}

	p := &prompt.Prompt{
		Name:    req.Name,
		Type:    req.Type,
		Content: req.Content,
	}

	err := c.promptService.CreatePrompt(p)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}

	r.Response.WriteJsonExit(g.Map{
		"code": 200,
		"msg":  "创建成功",
		"data": p,
	})
}

// GetPrompt 获取提示词信息
func (c *PromptController) GetPrompt(r *ghttp.Request) {
	id := r.Get("id").Uint64()
	p, err := c.promptService.GetPrompt(id)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}

	r.Response.WriteJsonExit(g.Map{
		"code": 200,
		"msg":  "获取成功",
		"data": p,
	})
}

// GetPromptsByType 获取指定类型的提示词列表
func (c *PromptController) GetPromptsByType(r *ghttp.Request) {
	promptType := r.Get("type").Int()
	prompts, err := c.promptService.GetPromptsByType(promptType)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}

	r.Response.WriteJsonExit(g.Map{
		"code": 200,
		"msg":  "获取成功",
		"data": prompts,
	})
}

// UpdatePrompt 更新提示词
func (c *PromptController) UpdatePrompt(r *ghttp.Request) {
	var req struct {
		ID      uint64 `json:"id" v:"required"`
		Name    string `json:"name"`
		Type    int    `json:"type" v:"in:1,2"`
		Content string `json:"content"`
	}

	if err := r.Parse(&req); err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 400,
			"msg":  err.Error(),
		})
	}

	p, err := c.promptService.GetPrompt(req.ID)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}

	if req.Name != "" {
		p.Name = req.Name
	}
	if req.Type != 0 {
		p.Type = req.Type
	}
	if req.Content != "" {
		p.Content = req.Content
	}

	err = c.promptService.UpdatePrompt(p)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 500,
			"msg":  err.Error(),
		})
	}

	r.Response.WriteJsonExit(g.Map{
		"code": 200,
		"msg":  "更新成功",
		"data": p,
	})
}

// DeletePrompt 删除提示词
func (c *PromptController) DeletePrompt(r *ghttp.Request) {
	id := r.Get("id").Uint64()
	err := c.promptService.DeletePrompt(id)
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
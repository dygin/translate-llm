package ai

import (
	"context"
)

// ModelType AI模型类型
type ModelType string

const (
	ModelTypeGemini ModelType = "gemini"
)

// ModelConfig AI模型配置
type ModelConfig struct {
	Type      ModelType `json:"type"`
	APIKey    string    `json:"api_key"`
	ModelName string    `json:"model_name"`
}

// GenerateContentRequest 生成内容请求
type GenerateContentRequest struct {
	VideoURL    string `json:"video_url"`
	SubtitleURL string `json:"subtitle_url"`
	Prompt      string `json:"prompt"`
}

// GenerateContentResponse 生成内容响应
type GenerateContentResponse struct {
	Content string `json:"content"`
}

// TranslationRequest 翻译请求
type TranslationRequest struct {
	Content        string `json:"content"`
	TargetLanguage string `json:"target_language"`
	Terminology    string `json:"terminology"`
	Prompt         string `json:"prompt"`
}

// TranslationResponse 翻译响应
type TranslationResponse struct {
	TranslatedContent string `json:"translated_content"`
}

// ModelService AI模型服务接口
type ModelService interface {
	GenerateContent(ctx context.Context, req *GenerateContentRequest) (*GenerateContentResponse, error)
	Translate(ctx context.Context, req *TranslationRequest) (*TranslationResponse, error)
}

// ModelFactory AI模型工厂接口
type ModelFactory interface {
	CreateModel(config *ModelConfig) (ModelService, error)
} 
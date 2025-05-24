package ai

import (
	"context"
	"ai-translate/internal/domain/ai"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type geminiService struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

// NewGeminiService 创建Gemini AI模型服务实例
func NewGeminiService(config *ai.ModelConfig) (ai.ModelService, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(config.APIKey))
	if err != nil {
		return nil, err
	}

	model := client.GenerativeModel(config.ModelName)
	
	return &geminiService{
		client: client,
		model:  model,
	}, nil
}

func (s *geminiService) GenerateContent(ctx context.Context, req *ai.GenerateContentRequest) (*ai.GenerateContentResponse, error) {
	// 构建提示词
	prompt := req.Prompt + "\n视频URL: " + req.VideoURL + "\n字幕URL: " + req.SubtitleURL
	
	// 生成内容
	resp, err := s.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}

	// 获取生成的内容
	content := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			content += string(text)
		}
	}

	return &ai.GenerateContentResponse{
		Content: content,
	}, nil
}

func (s *geminiService) Translate(ctx context.Context, req *ai.TranslationRequest) (*ai.TranslationResponse, error) {
	// 构建提示词
	prompt := req.Prompt + "\n目标语言: " + req.TargetLanguage + "\n术语表: " + req.Terminology + "\n待翻译内容: " + req.Content
	
	// 生成翻译
	resp, err := s.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}

	// 获取翻译结果
	translatedContent := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			translatedContent += string(text)
		}
	}

	return &ai.TranslationResponse{
		TranslatedContent: translatedContent,
	}, nil
}

type geminiFactory struct{}

// NewGeminiFactory 创建Gemini AI模型工厂实例
func NewGeminiFactory() ai.ModelFactory {
	return &geminiFactory{}
}

func (f *geminiFactory) CreateModel(config *ai.ModelConfig) (ai.ModelService, error) {
	return NewGeminiService(config)
} 
package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OpenAIService OpenAI服务实现
type OpenAIService struct {
	BaseAIService
	apiKey      string
	model       string
	maxTokens   int
	temperature float64
	timeout     int
}

// NewOpenAIService 创建OpenAI服务
func NewOpenAIService(apiKey, model string, maxTokens int, temperature float64, timeout int) (*OpenAIService, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API密钥不能为空")
	}

	service := &OpenAIService{
		apiKey:      apiKey,
		model:       model,
		maxTokens:   maxTokens,
		temperature: temperature,
		timeout:     timeout,
	}
	service.driver = service
	return service, nil
}

// Generate 生成内容
func (s *OpenAIService) Generate(ctx context.Context, prompt string) (string, error) {
	// 构建请求
	req := OpenAIRequest{
		Model:       s.model,
		MaxTokens:   s.maxTokens,
		Temperature: s.temperature,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	// 发送请求
	resp, err := s.sendRequest(ctx, "https://api.openai.com/v1/chat/completions", req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}

	// 解析响应
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("未获取到生成内容")
	}

	return resp.Choices[0].Message.Content, nil
}

// OpenAIRequest OpenAI请求结构
type OpenAIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
}

// Message 消息结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse OpenAI响应结构
type OpenAIResponse struct {
	Choices []Choice `json:"choices"`
}

// Choice 选择结构
type Choice struct {
	Message Message `json:"message"`
}

// GenerateContent 生成内容
func (s *OpenAIService) GenerateContent(ctx context.Context, prompt string, language string) (string, error) {
	// 构建提示词
	systemPrompt := fmt.Sprintf("你是一个专业的内容生成助手。请用%s语言生成内容。", language)
	userPrompt := fmt.Sprintf("请根据以下内容生成相关内容：\n%s", prompt)

	// 构建请求
	req := OpenAIRequest{
		Model:       s.model,
		MaxTokens:   s.maxTokens,
		Temperature: s.temperature,
		Messages: []Message{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: userPrompt,
			},
		},
	}

	// 发送请求
	resp, err := s.sendRequest(ctx, "https://api.openai.com/v1/chat/completions", req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}

	// 解析响应
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("未获取到生成内容")
	}

	return resp.Choices[0].Message.Content, nil
}

// Translate 翻译内容
func (s *OpenAIService) Translate(ctx context.Context, content string, sourceLang string, targetLang string) (string, error) {
	// 构建提示词
	systemPrompt := fmt.Sprintf("你是一个专业的翻译助手。请将%s语言翻译成%s语言。", sourceLang, targetLang)
	userPrompt := fmt.Sprintf("请翻译以下内容：\n%s", content)

	// 构建请求
	req := OpenAIRequest{
		Model:       s.model,
		MaxTokens:   s.maxTokens,
		Temperature: s.temperature,
		Messages: []Message{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: userPrompt,
			},
		},
	}

	// 发送请求
	resp, err := s.sendRequest(ctx, "https://api.openai.com/v1/chat/completions", req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}

	// 解析响应
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("未获取到翻译结果")
	}

	return resp.Choices[0].Message.Content, nil
}

// sendRequest 发送HTTP请求
func (s *OpenAIService) sendRequest(ctx context.Context, url string, req OpenAIRequest) (*OpenAIResponse, error) {
	// 序列化请求
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: time.Duration(s.timeout) * time.Second,
	}

	// 发送请求
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码：%d，响应：%s", resp.StatusCode, string(body))
	}

	// 解析响应
	var openAIResp OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &openAIResp, nil
} 
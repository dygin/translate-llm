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

// GeminiService Gemini服务实现
type GeminiService struct {
	BaseAIService
	apiKey  string
	model   string
	timeout int
}

// NewGeminiService 创建Gemini服务
func NewGeminiService(apiKey, model string, timeout int) (*GeminiService, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API密钥不能为空")
	}

	service := &GeminiService{
		apiKey:  apiKey,
		model:   model,
		timeout: timeout,
	}
	service.driver = service
	return service, nil
}

// Generate 生成内容
func (s *GeminiService) Generate(ctx context.Context, prompt string) (string, error) {
	// 构建请求
	req := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{
						Text: prompt,
					},
				},
			},
		},
	}

	// 发送请求
	resp, err := s.sendRequest(ctx, fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", s.model, s.apiKey), req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}

	// 解析响应
	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("未获取到生成内容")
	}

	return resp.Candidates[0].Content.Parts[0].Text, nil
}

// GeminiRequest Gemini请求结构
type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

// GeminiContent Gemini内容结构
type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

// GeminiPart Gemini部分结构
type GeminiPart struct {
	Text string `json:"text"`
}

// GeminiResponse Gemini响应结构
type GeminiResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
}

// GeminiCandidate Gemini候选结构
type GeminiCandidate struct {
	Content GeminiContent `json:"content"`
}

// sendRequest 发送HTTP请求
func (s *GeminiService) sendRequest(ctx context.Context, url string, req GeminiRequest) (*GeminiResponse, error) {
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
	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &geminiResp, nil
} 
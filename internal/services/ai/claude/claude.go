package claude

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"ai-assistant/internal/app/config"
	"ai-assistant/pkg/logger"
)

type ClaudeService struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	logger     *logger.Logger
}

type ClaudeRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ClaudeResponse struct {
	Content []Content `json:"content"`
	Error   *APIError `json:"error,omitempty"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type APIError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func NewClaudeService(cfg *config.Config) *ClaudeService {
	if cfg.AI.ClaudeAPIKey == "" {
		return nil // Claude is optional
	}

	return &ClaudeService{
		apiKey:     cfg.AI.ClaudeAPIKey,
		baseURL:    "https://api.anthropic.com/v1",
		httpClient: &http.Client{},
		logger:     logger.New(),
	}
}

func (c *ClaudeService) GenerateResponse(prompt string) (string, error) {
	c.logger.Infof("Generating Claude response for prompt: %s", prompt[:min(50, len(prompt))]+"...")

	reqBody := ClaudeRequest{
		Model:     "claude-3-haiku-20240307",
		MaxTokens: 2048,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Errorf("Claude API error: %s", string(body))
		return "", fmt.Errorf("claude API error: %s", resp.Status)
	}

	var claudeResp ClaudeResponse
	if err := json.Unmarshal(body, &claudeResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if claudeResp.Error != nil {
		return "", fmt.Errorf("claude API error: %s", claudeResp.Error.Message)
	}

	if len(claudeResp.Content) == 0 {
		return "", fmt.Errorf("no content returned from Claude")
	}

	return claudeResp.Content[0].Text, nil
}

func (c *ClaudeService) Close() error {
	// No resources to close for HTTP client
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
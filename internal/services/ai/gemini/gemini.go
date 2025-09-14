package gemini

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"ai-assistant/internal/app/config"
	"ai-assistant/pkg/logger"
)

type GeminiService struct {
	client *genai.Client
	model  *genai.GenerativeModel
	ctx    context.Context
	logger *logger.Logger
}

func NewGeminiService(cfg *config.Config) (*GeminiService, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.AI.GeminiAPIKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	model := client.GenerativeModel("gemini-1.5-flash")
	model.SetTemperature(0.7)
	model.SetTopK(32)
	model.SetTopP(0.9)
	model.SetMaxOutputTokens(2048)

	return &GeminiService{
		client: client,
		model:  model,
		ctx:    ctx,
		logger: logger.New(),
	}, nil
}

func (g *GeminiService) GenerateResponse(prompt string) (string, error) {
	g.logger.Infof("Generating response for prompt: %s", prompt[:min(50, len(prompt))]+"...")

	resp, err := g.model.GenerateContent(g.ctx, genai.Text(prompt))
	if err != nil {
		g.logger.Errorf("Failed to generate content: %v", err)
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no candidates returned from Gemini")
	}

	candidate := resp.Candidates[0]
	if len(candidate.Content.Parts) == 0 {
		return "", fmt.Errorf("no content parts returned from Gemini")
	}

	if textPart, ok := candidate.Content.Parts[0].(genai.Text); ok {
		return string(textPart), nil
	}

	return "", fmt.Errorf("unexpected content type from Gemini")
}

func (g *GeminiService) Close() error {
	return g.client.Close()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
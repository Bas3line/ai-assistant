package usecase

import (
	"context"

	"ai-assistant/internal/models"
	"ai-assistant/internal/repository"
	"ai-assistant/pkg/cache"
	"ai-assistant/pkg/errors"
)

// AIProvider interface for AI services
type AIProvider interface {
	GenerateResponse(prompt string) (string, error)
	Close() error
}

type AIUsecase struct {
	geminiService AIProvider
	claudeService AIProvider
	redisService  *cache.RedisService
}

func NewAIUsecase(geminiService AIProvider, claudeService AIProvider, redisService *cache.RedisService) *AIUsecase {
	return &AIUsecase{
		geminiService: geminiService,
		claudeService: claudeService,
		redisService:  redisService,
	}
}

func (u *AIUsecase) ProcessAIRequest(ctx context.Context, userID string, req *models.AIRequest) (*models.AIResponse, error) {
	if req.Prompt == "" {
		return nil, errors.ErrBadRequest("Prompt is required")
	}

	if req.Provider == "" {
		req.Provider = "gemini"
	}

	var response string
	var provider string
	var err error

	switch req.Provider {
	case "claude":
		if u.claudeService == nil {
			return nil, errors.ErrServiceUnavailable("Claude service not available")
		}
		response, err = u.claudeService.GenerateResponse(req.Prompt)
		provider = "claude"
	case "gemini":
		response, err = u.geminiService.GenerateResponse(req.Prompt)
		provider = "gemini"
	default:
		return nil, errors.ErrBadRequest("Invalid provider. Use 'gemini' or 'claude'")
	}

	if err != nil {
		return nil, errors.ErrInternalServerError("Failed to generate response: " + err.Error())
	}

	return &models.AIResponse{
		Response: response,
		Provider: provider,
	}, nil
}

type AuthUsecase struct {
	userRepo *repository.UserRepository
}

func NewAuthUsecase(userRepo *repository.UserRepository) *AuthUsecase {
	return &AuthUsecase{
		userRepo: userRepo,
	}
}

func (u *AuthUsecase) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.ErrDatabaseError
	}
	if user == nil {
		return nil, errors.ErrUserNotFound
	}
	return user, nil
}
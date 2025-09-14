package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"ai-assistant/internal/handlers"
	"ai-assistant/internal/models"
	"ai-assistant/internal/services/auth"
)

type MockAIUsecase struct {
	mock.Mock
}

func (m *MockAIUsecase) ProcessAIRequest(ctx context.Context, userID string, req *models.AIRequest) (*models.AIResponse, error) {
	args := m.Called(ctx, userID, req)
	return args.Get(0).(*models.AIResponse), args.Error(1)
}

func mockAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := &models.AuthUser{
			ID:    "user123",
			Email: "test@example.com",
		}
		ctx := context.WithValue(r.Context(), auth.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func TestAIHandler_Ask(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    models.AIRequest
		expectedStatus int
		expectedResp   models.AIResponse
		setupMock     func(*MockAIUsecase)
	}{
		{
			name: "successful gemini request",
			requestBody: models.AIRequest{
				Prompt:   "Hello, how are you?",
				Provider: "gemini",
			},
			expectedStatus: http.StatusOK,
			expectedResp: models.AIResponse{
				Response: "Hello! I'm doing well, thank you for asking.",
				Provider: "gemini",
			},
			setupMock: func(m *MockAIUsecase) {
				m.On("ProcessAIRequest", mock.Anything, "user123", mock.AnythingOfType("*models.AIRequest")).
					Return(&models.AIResponse{
						Response: "Hello! I'm doing well, thank you for asking.",
						Provider: "gemini",
					}, nil)
			},
		},
		{
			name: "successful claude request",
			requestBody: models.AIRequest{
				Prompt:   "What is Go programming language?",
				Provider: "claude",
			},
			expectedStatus: http.StatusOK,
			expectedResp: models.AIResponse{
				Response: "Go is a statically typed, compiled programming language.",
				Provider: "claude",
			},
			setupMock: func(m *MockAIUsecase) {
				m.On("ProcessAIRequest", mock.Anything, "user123", mock.AnythingOfType("*models.AIRequest")).
					Return(&models.AIResponse{
						Response: "Go is a statically typed, compiled programming language.",
						Provider: "claude",
					}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockAIUsecase)
			tt.setupMock(mockUsecase)

			handler := handlers.NewAIHandler(mockUsecase)

			router := chi.NewRouter()
			
			router.Use(mockAuthMiddleware)

			router.Route("/ai", func(r chi.Router) {
				handler.RegisterRoutes(r)
			})

			requestBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/ai/ask", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if w.Code == http.StatusOK {
				var response models.AIResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResp.Response, response.Response)
				assert.Equal(t, tt.expectedResp.Provider, response.Provider)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}
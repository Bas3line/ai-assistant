package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"ai-assistant/internal/models"
	"ai-assistant/internal/services/auth"
)

// AIUsecaseInterface defines the interface for AI usecase
type AIUsecaseInterface interface {
	ProcessAIRequest(ctx context.Context, userID string, req *models.AIRequest) (*models.AIResponse, error)
}

type AIHandler struct {
	aiUsecase AIUsecaseInterface
}

func NewAIHandler(aiUsecase AIUsecaseInterface) *AIHandler {
	return &AIHandler{
		aiUsecase: aiUsecase,
	}
}

func (h *AIHandler) Ask(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetCurrentUser(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		response := map[string]interface{}{"error": "Unauthorized"}
		json.NewEncoder(w).Encode(response)
		return
	}

	var req models.AIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"error": "Invalid JSON"}
		json.NewEncoder(w).Encode(response)
		return
	}

	response, err := h.aiUsecase.ProcessAIRequest(r.Context(), user.ID, &req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := map[string]interface{}{"error": err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AIHandler) RegisterRoutes(router chi.Router) {
	router.Post("/ask", h.Ask)
}
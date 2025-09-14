package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"ai-assistant/internal/services/auth"
)

type EmailHandler struct {
}

func NewEmailHandler() *EmailHandler {
	return &EmailHandler{}
}

func (h *EmailHandler) GetEmails(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetCurrentUser(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		response := map[string]interface{}{"error": "Unauthorized"}
		json.NewEncoder(w).Encode(response)
		return
	}

	_ = user

	response := map[string]interface{}{
		"message": "Email fetching not implemented yet",
		"user_id": user.ID,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(response)
}

func (h *EmailHandler) SendEmail(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetCurrentUser(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		response := map[string]interface{}{"error": "Unauthorized"}
		json.NewEncoder(w).Encode(response)
		return
	}

	_ = user

	response := map[string]interface{}{
		"message": "Email sending not implemented yet",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(response)
}

func (h *EmailHandler) RegisterRoutes(router chi.Router) {
	router.Get("/", h.GetEmails)
	router.Post("/send", h.SendEmail)
}
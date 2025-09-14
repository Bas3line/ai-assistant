package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	"github.com/go-chi/chi/v5"
	"ai-assistant/internal/models"
	"ai-assistant/internal/services/auth"
)

type AuthHandler struct {
	authService *auth.AuthService
}

func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) GoogleOAuth(w http.ResponseWriter, r *http.Request) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	if clientID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"error": "Google OAuth not configured"}
		json.NewEncoder(w).Encode(response)
		return
	}

	state := generateRandomState()
	baseURL := os.Getenv("BASE_URL")
	redirectURI := baseURL + "/api/auth/google/callback"

	params := url.Values{}
	params.Add("client_id", clientID)
	params.Add("redirect_uri", redirectURI)
	params.Add("response_type", "code")
	params.Add("scope", "openid email profile")
	params.Add("state", state)

	authorizeURL := "https://accounts.google.com/o/oauth2/v2/auth?" + params.Encode()

	response := map[string]interface{}{
		"authorization_url": authorizeURL,
		"state": state,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	errorParam := r.URL.Query().Get("error")

	if errorParam != "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"error": "OAuth error: " + errorParam}
		json.NewEncoder(w).Encode(response)
		return
	}

	if code == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"error": "No authorization code provided"}
		json.NewEncoder(w).Encode(response)
		return
	}

	mockUser := &models.User{
		ID:    "demo_user_123",
		Email: "demo@example.com",
		Name:  stringPtr("Demo User"),
		Image: stringPtr("https://example.com/avatar.jpg"),
	}

	token, err := h.authService.GenerateToken(mockUser)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"error": "Failed to generate token"}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"message": "Authentication successful",
		"token":   token,
		"user": map[string]interface{}{
			"id":    mockUser.ID,
			"email": mockUser.Email,
			"name":  mockUser.Name,
			"image": mockUser.Image,
		},
		"state": state,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetCurrentUser(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		response := map[string]interface{}{"error": "Unauthorized"}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{"message": "Logged out successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func generateRandomState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func stringPtr(s string) *string {
	return &s
}

func (h *AuthHandler) RegisterRoutes(router chi.Router) {
	router.Route("/auth", func(r chi.Router) {
		r.Get("/google", h.GoogleOAuth)
		r.Get("/google/callback", h.GoogleCallback)
		r.With(h.authService.RequireAuth()).Get("/me", h.Me)
		r.With(h.authService.RequireAuth()).Post("/logout", h.Logout)
	})
}
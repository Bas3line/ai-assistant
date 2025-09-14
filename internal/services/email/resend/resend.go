package resend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"ai-assistant/internal/app/config"
	"ai-assistant/pkg/logger"
)

type ResendService struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	logger     *logger.Logger
}

type EmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html,omitempty"`
	Text    string   `json:"text,omitempty"`
}

type EmailResponse struct {
	ID string `json:"id"`
}

func NewResendService(cfg *config.Config) *ResendService {
	return &ResendService{
		apiKey:     cfg.Email.ResendAPIKey,
		baseURL:    "https://api.resend.com",
		httpClient: &http.Client{},
		logger:     logger.New(),
	}
}

func (r *ResendService) SendEmail(req EmailRequest) (*EmailResponse, error) {
	r.logger.Infof("Sending email to %v", req.To)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", r.baseURL+"/emails", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+r.apiKey)

	resp, err := r.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("email API error: %s", resp.Status)
	}

	var emailResp EmailResponse
	if err := json.NewDecoder(resp.Body).Decode(&emailResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	r.logger.Infof("Email sent successfully with ID: %s", emailResp.ID)
	return &emailResp, nil
}
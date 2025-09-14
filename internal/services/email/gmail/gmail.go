package gmail

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"ai-assistant/internal/app/config"
	"ai-assistant/pkg/logger"
)

type GmailService struct {
	service *gmail.Service
	logger  *logger.Logger
}

type GmailMessage struct {
	ID       string
	ThreadID string
	Subject  string
	From     string
	To       []string
	Body     string
	HTMLBody string
	Labels   []string
	IsRead   bool
}

func NewGmailService(cfg *config.Config, token *oauth2.Token) (*GmailService, error) {
	ctx := context.Background()
	
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.Google.ClientID,
		ClientSecret: cfg.Google.ClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/gmail.readonly",
			"https://www.googleapis.com/auth/gmail.send",
		},
	}

	client := oauthConfig.Client(ctx, token)
	service, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gmail service: %w", err)
	}

	return &GmailService{
		service: service,
		logger:  logger.New(),
	}, nil
}

func (g *GmailService) GetMessages(userID string, maxResults int64) ([]*GmailMessage, error) {
	call := g.service.Users.Messages.List(userID).MaxResults(maxResults)
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}

	var messages []*GmailMessage
	for _, message := range response.Messages {
		msg, err := g.GetMessage(userID, message.Id)
		if err != nil {
			g.logger.Errorf("Failed to get message %s: %v", message.Id, err)
			continue
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (g *GmailService) GetMessage(userID, messageID string) (*GmailMessage, error) {
	message, err := g.service.Users.Messages.Get(userID, messageID).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	gmailMsg := &GmailMessage{
		ID:       message.Id,
		ThreadID: message.ThreadId,
		Labels:   message.LabelIds,
		IsRead:   !contains(message.LabelIds, "UNREAD"),
	}

	for _, header := range message.Payload.Headers {
		switch header.Name {
		case "Subject":
			gmailMsg.Subject = header.Value
		case "From":
			gmailMsg.From = header.Value
		case "To":
			gmailMsg.To = strings.Split(header.Value, ",")
		}
	}

	body, htmlBody := extractBody(message.Payload)
	gmailMsg.Body = body
	gmailMsg.HTMLBody = htmlBody

	return gmailMsg, nil
}

func (g *GmailService) SendMessage(userID string, to, subject, body string) error {
	message := fmt.Sprintf("To: %s\nSubject: %s\n\n%s", to, subject, body)
	encodedMessage := base64.URLEncoding.EncodeToString([]byte(message))

	gmailMessage := &gmail.Message{
		Raw: encodedMessage,
	}

	_, err := g.service.Users.Messages.Send(userID, gmailMessage).Do()
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	g.logger.Infof("Message sent successfully to %s", to)
	return nil
}

func (g *GmailService) MarkAsRead(userID, messageID string) error {
	req := &gmail.ModifyMessageRequest{
		RemoveLabelIds: []string{"UNREAD"},
	}

	_, err := g.service.Users.Messages.Modify(userID, messageID, req).Do()
	if err != nil {
		return fmt.Errorf("failed to mark message as read: %w", err)
	}

	return nil
}

func extractBody(payload *gmail.MessagePart) (string, string) {
	var body, htmlBody string

	if payload.Body != nil && payload.Body.Data != "" {
		decoded, err := base64.URLEncoding.DecodeString(payload.Body.Data)
		if err == nil {
			if payload.MimeType == "text/html" {
				htmlBody = string(decoded)
			} else {
				body = string(decoded)
			}
		}
	}

	for _, part := range payload.Parts {
		if part.MimeType == "text/plain" && part.Body != nil && part.Body.Data != "" {
			decoded, err := base64.URLEncoding.DecodeString(part.Body.Data)
			if err == nil {
				body = string(decoded)
			}
		} else if part.MimeType == "text/html" && part.Body != nil && part.Body.Data != "" {
			decoded, err := base64.URLEncoding.DecodeString(part.Body.Data)
			if err == nil {
				htmlBody = string(decoded)
			}
		}
	}

	return body, htmlBody
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
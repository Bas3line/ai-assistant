package usecase

import (
	"context"
	"fmt"

	"ai-assistant/internal/models"
	"ai-assistant/internal/repository"
	"ai-assistant/internal/services/email/gmail"
	"ai-assistant/internal/services/email/resend"
	"ai-assistant/pkg/errors"
)

type EmailUsecase struct {
	emailRepo    *repository.EmailRepository
	resendSvc    *resend.ResendService
	gmailSvc     *gmail.GmailService
}

func NewEmailUsecase(emailRepo *repository.EmailRepository, resendSvc *resend.ResendService, gmailSvc *gmail.GmailService) *EmailUsecase {
	return &EmailUsecase{
		emailRepo: emailRepo,
		resendSvc: resendSvc,
		gmailSvc:  gmailSvc,
	}
}

func (u *EmailUsecase) GetUserEmails(ctx context.Context, userID string, limit, offset int) ([]*models.Email, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	emails, err := u.emailRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, errors.ErrDatabaseError
	}
	
	return emails, nil
}

func (u *EmailUsecase) SendEmail(ctx context.Context, from string, to []string, subject, body string) error {
	if u.resendSvc == nil {
		return errors.ErrServiceUnavailable("Email service not configured")
	}
	
	req := resend.EmailRequest{
		From:    from,
		To:      to,
		Subject: subject,
		Text:    body,
	}
	
	_, err := u.resendSvc.SendEmail(req)
	if err != nil {
		return errors.ErrInternalServerError(fmt.Sprintf("Failed to send email: %v", err))
	}
	
	return nil
}

func (u *EmailUsecase) SyncGmailEmails(ctx context.Context, userID string) error {
	if u.gmailSvc == nil {
		return errors.ErrServiceUnavailable("Gmail service not configured")
	}
	
	messages, err := u.gmailSvc.GetMessages(userID, 50)
	if err != nil {
		return errors.ErrExternalService
	}
	
	for _, msg := range messages {
		email := &models.Email{
			ID:        generateID(),
			MessageID: msg.ID,
			ThreadID:  &msg.ThreadID,
			Subject:   &msg.Subject,
			From:      msg.From,
			To:        msg.To,
			Body:      &msg.Body,
			HTMLBody:  &msg.HTMLBody,
			IsRead:    msg.IsRead,
			Labels:    msg.Labels,
			UserID:    userID,
		}
		
		if err := u.emailRepo.Create(ctx, email); err != nil {
			continue
		}
	}
	
	return nil
}

func (u *EmailUsecase) MarkEmailAsRead(ctx context.Context, emailID string) error {
	return u.emailRepo.UpdateReadStatus(ctx, emailID, true)
}

func generateID() string {
	return fmt.Sprintf("email_%d", len("placeholder"))
}
package utils

import (
	"fmt"
	"html"

	"github.com/dariubs/scaffold/app/config"
	"github.com/resend/resend-go/v3"
)

type EmailService struct {
	client *resend.Client
	from   string
}

// NewEmailService creates an email service using Resend. If RESEND_API_KEY or
// RESEND_FROM is not set, returns (nil, nil) so callers can skip sending.
func NewEmailService() (*EmailService, error) {
	if config.C.Resend.APIKey == "" || config.C.Resend.From == "" {
		return nil, nil
	}
	client := resend.NewClient(config.C.Resend.APIKey)
	return &EmailService{client: client, from: config.C.Resend.From}, nil
}

// SendWelcome sends a welcome email to the given address. userName may be empty.
func (s *EmailService) SendWelcome(toEmail, userName string) error {
	greeting := "Welcome to Scaffold!"
	if userName != "" {
		greeting = fmt.Sprintf("Hi %s, welcome to Scaffold!", html.EscapeString(userName))
	}
	body := fmt.Sprintf(`<p>%s</p><p>Thanks for signing up.</p>`, greeting)

	_, err := s.client.Emails.Send(&resend.SendEmailRequest{
		From:    s.from,
		To:      []string{toEmail},
		Subject: "Welcome to Scaffold",
		Html:    body,
	})
	if err != nil {
		Logger.Error("Failed to send welcome email", "err", err, "to", toEmail)
		return err
	}
	return nil
}

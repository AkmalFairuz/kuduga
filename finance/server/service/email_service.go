package service

import "github.com/akmalfairuz/finance/server/provider"

type EmailService struct {
	emailProvider provider.EmailProvider
}

func NewEmailService(emailProvider provider.EmailProvider) *EmailService {
	return &EmailService{
		emailProvider: emailProvider,
	}
}

func (s *EmailService) SendEmail(to []string, subject, body, from string) error {
	return s.emailProvider.SendEmail(to, subject, body, from)
}

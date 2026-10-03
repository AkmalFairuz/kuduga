package provider

type EmailProvider interface {
	SendEmail(to []string, subject, body, from string) error
}

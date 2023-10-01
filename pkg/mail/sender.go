package mail

import (
	"fmt"
	"net/smtp"
	"os"
)

var Sender *sender

func Initialize() {
	var (
		host     = os.Getenv("SMTP_HOST")
		port     = os.Getenv("SMTP_PORT")
		email    = os.Getenv("SMTP_EMAIL")
		password = os.Getenv("SMTP_PASSWORD")
	)

	Sender = &sender{
		Auth:  smtp.PlainAuth("", email, password, host),
		Email: email,
		Addr:  fmt.Sprintf("%s:%s", host, port),
	}
}

type sender struct {
	Auth  smtp.Auth
	Email string
	Addr  string
}

func (s sender) Send(subject, message string, to []string) error {
	return smtp.SendMail(s.Addr, s.Auth, s.Email, to, formatMessage(subject, message))
}

func formatMessage(subject, message string) []byte {
	return []byte(fmt.Sprintf("Subject: %s\n%s", subject, message))
}

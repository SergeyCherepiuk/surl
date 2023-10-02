package mail

import (
	"fmt"
	"net/smtp"
	"os"
)

var Sender sender

type sender struct {
	Auth     smtp.Auth
	Email    string
	SmtpAddr string
}

func Initialize() {
	var (
		host     = os.Getenv("SMTP_HOST")
		port     = os.Getenv("SMTP_PORT")
		email    = os.Getenv("SMTP_EMAIL")
		password = os.Getenv("SMTP_PASSWORD")
	)

	Sender = sender{
		Auth:     smtp.PlainAuth("", email, password, host),
		Email:    email,
		SmtpAddr: fmt.Sprintf("%s:%s", host, port),
	}
}

func (s sender) Send(to, subject, message string) error {
	return smtp.SendMail(
		s.SmtpAddr, s.Auth, s.Email, []string{to},
		formatMessage(s.Email, to, subject, message),
	)
}

func formatMessage(from, to, subject, message string) []byte {
	return []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		from, to, subject, message,
	))
}

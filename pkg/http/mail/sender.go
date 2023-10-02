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

func (s sender) SendMessage(to, subject, message string) error {
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

func (s sender) SendHTML(to, subject, html string) error {
	return smtp.SendMail(
		s.SmtpAddr, s.Auth, s.Email, []string{to},
		formatHTML(s.Email, to, subject, html),
	)
}

func formatHTML(from, to, subject, html string) []byte {
	return []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s",
		from, to, subject, html,
	))
}

package mail

import (
	"bytes"
	"fmt"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/SergeyCherepiuk/surl/pkg/http/template"
	"github.com/SergeyCherepiuk/surl/public/views/pages"
)

type verificationSender struct{}

func NewVerificationSender() *verificationSender {
	return &verificationSender{}
}

func (vs verificationSender) Send(email string, verificationRequest domain.VerificationRequest) error {
	var html bytes.Buffer

	link := fmt.Sprintf(
		"http://localhost:3000/api/verification/%s/%s",
		verificationRequest.Username, verificationRequest.ID,
	)
	data := pages.MailPageData{
		Username: verificationRequest.Username,
		Paragraphs: []string{
			"Please click the button below to verify your account.",
			"This verification link will expire in 48 hours. You always can request a new one.",
		},
		Link:       link,
		ButtonText: "Verify",
	}
	if err := template.Renderer.Templates["mail"].Execute(&html, data); err != nil {
		return err
	}

	return Sender.SendHTML(email, "Account verification", html.String())
}

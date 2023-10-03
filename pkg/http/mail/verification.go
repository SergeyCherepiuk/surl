package mail

import (
	"bytes"
	"fmt"

	"github.com/SergeyCherepiuk/surl/pkg/http/template"
	"github.com/SergeyCherepiuk/surl/public/views/pages"
)

type verificationSender struct{}

func NewVerificationSender() *verificationSender {
	return &verificationSender{}
}

func (vs verificationSender) Send(email, username, id string) error {
	var html bytes.Buffer

	data := pages.VerificationMailPageData{
		Username: username,
		Link:     fmt.Sprintf("http://localhost:3000/api/verification/%s/%s", username, id),
	}
	if err := template.Renderer.Templates["verification-mail"].Execute(&html, data); err != nil {
		return err
	}

	return Sender.SendHTML(email, "Account verification", html.String())
}

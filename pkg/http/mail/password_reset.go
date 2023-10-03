package mail

import (
	"bytes"
	"context"
	"fmt"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/SergeyCherepiuk/surl/pkg/http/template"
	"github.com/SergeyCherepiuk/surl/public/views/pages"
)

type passwordResetSender struct{}

func NewPasswordResetSender() *passwordResetSender {
	return &passwordResetSender{}
}

func (prs passwordResetSender) Send(ctx context.Context, email string, passwordResetRequest domain.PasswordResetRequest) error {
	var html bytes.Buffer

	link := fmt.Sprintf(
		"http://localhost:3000/api/password-reset/%s/%s",
		passwordResetRequest.Username, passwordResetRequest.ID,
	)
	data := pages.MailPageData{
		Username:   passwordResetRequest.Username,
		Message:    "Please click the button below to provide a new password",
		Link:       link,
		ButtonText: "Reset password",
	}
	if err := template.Renderer.Templates["mail"].Execute(&html, data); err != nil {
		return err
	}

	return Sender.SendHTML(email, "Reset password", html.String())
}

package pages

import (
	"html/template"

	"github.com/SergeyCherepiuk/surl/public/views/components"
)

type PasswordResetPageData struct {
	Username                   string
	NewPasswordInputData       components.InputData
	NewPasswordRepeatInputData components.InputData
	ButtonData                 components.ButtonData
}

func PasswordResetPage() *template.Template {
	return template.Must(template.ParseFiles(
		"public/views/layout.html",
		"public/views/pages/password-reset.html",
		"public/views/components/input.html",
		"public/views/components/button.html",
	))
}

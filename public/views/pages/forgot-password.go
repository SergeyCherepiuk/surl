package pages

import (
	"html/template"

	"github.com/SergeyCherepiuk/surl/public/views/components"
)

type ForgotPasswordPageData struct {
	UsernameInputData components.InputData
	ButtonData        components.ButtonData
}

func ForgotPasswordPage() *template.Template {
	return template.Must(template.ParseFiles(
		"public/views/layout.html",
		"public/views/pages/forgot-password.html",
		"public/views/components/input.html",
		"public/views/components/button.html",
	))
}

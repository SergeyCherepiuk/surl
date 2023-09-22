package pages

import (
	"html/template"

	"github.com/SergeyCherepiuk/surl/public/views/components"
)

type LoginPageData struct {
	UsernameInputData components.InputData
	PasswordInputData components.InputData
	ButtonData        components.ButtonData
}

func LoginPage() *template.Template {
	return template.Must(template.ParseFiles(
		"public/views/layout.html",
		"public/views/pages/login.html",
		"public/views/components/input.html",
		"public/views/components/button.html",
	))
}

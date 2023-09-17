package pages

import (
	"html/template"

	"github.com/SergeyCherepiuk/surl/public/views/components"
)

type SignUpPageData struct {
	UsernameInputData components.InputComponentData
	PasswordInputData components.InputComponentData
	ButtonData        components.ButtonComponentData
}

func SignUpPage() *template.Template {
	return template.Must(template.ParseFiles(
		"public/views/layout.html",
		"public/views/pages/signup.html",
		"public/views/components/input.html",
		"public/views/components/button.html",
	))
}

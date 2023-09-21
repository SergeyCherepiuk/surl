package components

import (
	"html/template"
)

type DialogComponentData struct {
	Text                  string
	ConfirmIconButtonData IconButtonComponentData
	DeclineIconButtonData IconButtonComponentData
}

func DialogComponent() *template.Template {
	return template.Must(template.ParseFiles(
		"public/views/components/dialog.html",
		"public/views/components/icon-button.html",
	))
}

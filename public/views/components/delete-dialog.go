package components

import (
	"html/template"
)

type DeleteDialogComponentData struct {
	Text                  string
	ConfirmIconButtonData IconButtonComponentData
	DeclineIconButtonData IconButtonComponentData
}

func DeleteDialogComponent() *template.Template {
	return template.Must(template.ParseFiles(
		"public/views/components/delete-dialog.html",
		"public/views/components/icon-button.html",
	))
}

package components

import (
	"html/template"
)

type DeleteDialogData struct {
	Text                  string
	ConfirmIconButtonData IconButtonData
	DeclineIconButtonData IconButtonData
}

func DeleteDialog() *template.Template {
	return template.Must(template.ParseFiles(
		"public/views/components/delete-dialog.html",
		"public/views/components/icon-button.html",
	))
}

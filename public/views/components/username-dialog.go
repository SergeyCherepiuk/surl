package components

import "html/template"

type UsernameDialogComponentData struct {
	InputData             InputComponentData
	ConfirmIconButtonData IconButtonComponentData
	DeclineIconButtonData IconButtonComponentData
}

func UsernameDialogComponent() *template.Template {
	return template.Must(template.ParseFiles(
		"public/views/components/username-dialog.html",
		"public/views/components/input.html",
		"public/views/components/icon-button.html",
	))
}

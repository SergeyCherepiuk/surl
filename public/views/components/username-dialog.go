package components

import "html/template"

type UsernameDialogData struct {
	InputData             InputData
	ConfirmIconButtonData IconButtonData
	DeclineIconButtonData IconButtonData
}

func UsernameDialog() *template.Template {
	return template.Must(template.ParseFiles(
		"public/views/components/username-dialog.html",
		"public/views/components/input.html",
		"public/views/components/icon-button.html",
	))
}

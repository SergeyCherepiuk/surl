package components

import "html/template"

type PasswordDialogData struct {
	OldPasswordInputData       InputData
	NewPasswordInputData       InputData
	NewPasswordRepeatInputData InputData
	SubmitButtonData           ButtonData
	GoBackButtonData           ButtonData
}

func PasswordDialog() *template.Template {
	return template.Must(template.ParseFiles(
		"public/views/components/password-dialog.html",
		"public/views/components/input.html",
		"public/views/components/button.html",
	))
}
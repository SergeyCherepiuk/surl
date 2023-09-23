package components

import "html/template"

type OriginDialogData struct {
	Username         string
	Hash             string
	OriginInputData  InputData
	SubmitButtonData ButtonData
	GoBackButtonData ButtonData
}

func OriginDialog() *template.Template {
	return template.Must(template.ParseFiles(
		"public/views/components/origin-dialog.html",
		"public/views/components/input.html",
		"public/views/components/button.html",
	))
}

package components

import "html/template"

type IconsRowComponentData struct {
	ChangeUsernameIconButtonData IconButtonComponentData
	ChangePasswordIconButtonData IconButtonComponentData
	DeleteAccountIconButtonData  IconButtonComponentData
}

func IconsRowComponent() *template.Template {
	return template.Must(template.ParseFiles(
		"public/views/components/icons-row.html",
		"public/views/components/icon-button.html",
	))
}

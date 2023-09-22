package components

import "html/template"

type IconsRowData struct {
	ChangeUsernameIconButtonData IconButtonData
	ChangePasswordIconButtonData IconButtonData
	DeleteAccountIconButtonData  IconButtonData
}

func IconsRow() *template.Template {
	return template.Must(template.ParseFiles(
		"public/views/components/icons-row.html",
		"public/views/components/icon-button.html",
	))
}

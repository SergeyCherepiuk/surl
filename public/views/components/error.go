package components

import "html/template"

type ErrorData struct {
	Message string
}

func Error() *template.Template {
	return template.Must(template.ParseFiles(
		"public/views/components/error.html",
	))
}

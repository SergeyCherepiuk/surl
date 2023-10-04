package components

import "html/template"

func Error() *template.Template {
	return template.Must(template.ParseFiles(
		"public/views/components/error.html",
	))
}

package components

import "html/template"

func Info() *template.Template {
	return template.Must(template.ParseFiles(
		"public/views/components/info.html",
	))
}

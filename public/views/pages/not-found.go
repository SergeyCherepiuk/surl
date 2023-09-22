package pages

import "html/template"

func NotFoundPage() *template.Template {
	return template.Must(template.ParseFiles(
		"public/views/layout.html",
		"public/views/pages/not-found.html",
	))
}

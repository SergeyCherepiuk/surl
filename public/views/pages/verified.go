package pages

import "html/template"

func VerifiedPage() *template.Template {
	return template.Must(template.ParseFiles(
		"public/views/layout.html",
		"public/views/pages/verified.html",
	))
}

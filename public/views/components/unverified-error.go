package components

import "html/template"

func UnverifiedError() *template.Template {
	return template.Must(template.ParseFiles(
		"public/views/components/unverified-error.html",
	))
}

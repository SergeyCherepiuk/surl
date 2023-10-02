package pages

import "html/template"

type VerificationMailPageData struct {
	Username string
	Link     string
}

func VerificationMail() *template.Template {
	return template.Must(template.ParseFiles(
		"public/views/pages/verification-mail.html",
	))
}

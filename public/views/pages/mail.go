package pages

import "html/template"

type MailPageData struct {
	Username   string
	Message    string
	Link       string
	ButtonText string
}

func MailPage() *template.Template {
	return template.Must(template.ParseFiles(
		"public/views/pages/mail.html",
	))
}

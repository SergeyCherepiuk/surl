package pages

import (
	"html/template"
)

type HomePageData struct {
	Username string
}

func HomePage() *template.Template {
	return template.Must(template.ParseFiles(
		"public/views/layout.html",
		"public/views/pages/home.html",
		"public/views/components/icon-button.html",
		"public/views/components/url-input.html",
	))
}

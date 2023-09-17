package pages

import (
	"html/template"

	"github.com/SergeyCherepiuk/surl/public/views/components"
)

type HomePageData struct {
	Username      string
	UrlInputData  components.InputWithButtonComponentData
}

func HomePage() *template.Template {
	return template.Must(template.ParseFiles(
		"public/views/layout.html",
		"public/views/pages/home.html",
		"public/views/components/input-with-button.html",
	))
}

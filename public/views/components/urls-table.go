package components

import (
	"html/template"

	"github.com/SergeyCherepiuk/surl/domain"
)

type UrlsTableComponentData struct {
	Urls []domain.Url
}

func UrlsTableComponent() *template.Template {
	return template.Must(template.ParseFiles(
		"public/views/components/urls-table.html",
	))
}

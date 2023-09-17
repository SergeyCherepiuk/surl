package components

import (
	"html/template"
	"strings"

	"github.com/SergeyCherepiuk/surl/domain"
)

type UrlsTableComponentData struct {
	Urls []domain.Url
}

func UrlsTableComponent() *template.Template {
	return template.Must(template.New("urls-table.html").Funcs(template.FuncMap{
		"trimProtocol": func(url string) string {
			trimmed := strings.TrimPrefix(url, "http://")
			trimmed = strings.TrimPrefix(trimmed, "https://")
			return trimmed
		},
	}).ParseFiles(
		"public/views/components/urls-table.html",
	))
}

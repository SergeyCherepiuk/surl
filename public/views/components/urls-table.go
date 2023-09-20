package components

import (
	"html/template"
	"strings"
	"time"

	"github.com/SergeyCherepiuk/surl/domain"
)

type UrlsTableComponentData struct {
	Urls []domain.Url
}

func UrlsTableComponent() *template.Template {
	return template.Must(template.New("urls-table.html").Funcs(template.FuncMap{
		"formatLink": func(link string) string {
			l := strings.TrimPrefix(link, "http://")
			l = strings.TrimPrefix(l, "https://")
			return strings.TrimSuffix(l, "/")
		},
		"formatDate": func(date time.Time) string {
			return date.Format("02 Jan 2006")
		},
	}).ParseFiles(
		"public/views/components/urls-table.html",
	))
}

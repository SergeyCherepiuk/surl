package components

import (
	"html/template"
	"strings"
	"time"

	"github.com/SergeyCherepiuk/surl/domain"
)

type UrlsTableData struct {
	Urls     []domain.Url
	SortedBy string
	Reversed bool
}

func UrlsTable() *template.Template {
	return template.Must(template.New("urls-table.html").Funcs(template.FuncMap{
		"formatLink": func(link string) string {
			l := strings.TrimPrefix(link, "http://")
			l = strings.TrimPrefix(l, "https://")
			return strings.TrimSuffix(l, "/")
		},
		"formatDate": func(date time.Time) string {
			return date.Format("02 Jan 2006")
		},
		"formatDateWithTime": func(date time.Time) string {
			if date.Unix() == 0 {
				return "Never"
			}
			return date.Format("02 Jan 2006 15:04:05")
		},
	}).ParseFiles(
		"public/views/components/urls-table.html",
	))
}

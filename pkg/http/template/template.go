package template

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

var (
	Renderer = &Template{
		Templates: make(map[string]*template.Template),
	}
	notFoundTmpl = template.Must(template.ParseFiles(
		"public/views/layout.html",
		"public/views/404.html",
	))
)

func init() {
	// Register pages and components here
}

type Template struct {
	Templates map[string]*template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if tmpl, ok := t.Templates[name]; ok {
		return tmpl.Execute(w, data)
	}
	return notFoundTmpl.Execute(w, nil)
}

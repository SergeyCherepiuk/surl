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
	Renderer.Templates["greeting.html"] = template.Must(template.ParseFiles(
		"public/views/layout.html",
		"public/views/greeting.html",
	))
}

type Template struct {
	Templates map[string]*template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if tmpl, ok := t.Templates[name]; ok {
		return tmpl.ExecuteTemplate(w, "layout.html", data)
	}
	return notFoundTmpl.ExecuteTemplate(w, "layout.html", nil)
}

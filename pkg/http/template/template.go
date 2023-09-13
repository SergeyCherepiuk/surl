package template

import (
	"html/template"
	"io"
	"strings"

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
	// Pages
	Renderer.Templates["todo"] = template.Must(template.ParseFiles(
		"public/views/layout.html",
		"public/views/todo.html",
	))

	// Components
	Renderer.Templates["components/todo-list"] = template.Must(template.ParseFiles(
		"public/views/components/todo-list.html",
	))
}

type Template struct {
	Templates map[string]*template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.Templates[name]

	if ok && strings.HasPrefix(name, "components/") {
		return tmpl.Execute(w, data)
	}

	if ok {
		return tmpl.ExecuteTemplate(w, "layout.html", data)
	}

	return notFoundTmpl.ExecuteTemplate(w, "layout.html", nil)
}

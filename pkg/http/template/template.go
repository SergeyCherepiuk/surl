package template

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

var Renderer = &Template{
	Templates: make(map[string]*template.Template),
}

func init() {
	Renderer.Templates["home"] = template.Must(template.ParseFiles(
		"public/views/layout.html",
		"public/views/home.html",
	))
	Renderer.Templates["404"] = template.Must(template.ParseFiles(
		"public/views/layout.html",
		"public/views/404.html",
	))

	// Authentication
	Renderer.Templates["login"] = template.Must(template.ParseFiles(
		"public/views/layout.html",
		"public/views/auth/login.html",
	))
	Renderer.Templates["signup"] = template.Must(template.ParseFiles(
		"public/views/layout.html",
		"public/views/auth/signup.html",
	))

	// Components
	Renderer.Templates["components/urls-table-content"] = template.Must(template.ParseFiles(
		"public/views/components/urls-table-content.html",
	))
}

type Template struct {
	Templates map[string]*template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if tmpl, ok := t.Templates[name]; ok {
		return tmpl.Execute(w, data)
	}
	return t.Templates["404"].Execute(w, nil)
}

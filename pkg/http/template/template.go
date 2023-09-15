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
	Renderer.Templates["home"] = template.Must(template.ParseFiles(
		"public/views/layout.html",
		"public/views/home.html",
	))

	// Authentication
	Renderer.Templates["signup"] = template.Must(template.ParseFiles(
		"public/views/layout.html",
		"public/views/auth/signup.html",
	))

	// Components
	Renderer.Templates["component/error"] = template.Must(template.ParseFiles(
		"public/views/components/error.html",
	))
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

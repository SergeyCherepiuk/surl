package template

import (
	"html/template"
	"io"

	"github.com/SergeyCherepiuk/surl/public/views/components"
	"github.com/SergeyCherepiuk/surl/public/views/pages"
	"github.com/labstack/echo/v4"
)

var Renderer = &Template{
	Templates: make(map[string]*template.Template),
}

func init() {
	Renderer.Templates["home"] = pages.HomePage()
	Renderer.Templates["not-found"] = pages.NotFoundPage()
	Renderer.Templates["login"] = pages.LoginPage()
	Renderer.Templates["signup"] = pages.SignUpPage()

	Renderer.Templates["components/urls-table"] = components.UrlsTableComponent()
	Renderer.Templates["components/icons-row"] = components.IconsRowComponent()
	Renderer.Templates["components/dialog"] = components.DialogComponent()
}

type Template struct {
	Templates map[string]*template.Template
}

func (t *Template) Render(w io.Writer, name string, data any, c echo.Context) error {
	if tmpl, ok := t.Templates[name]; ok {
		return tmpl.Execute(w, data)
	}
	return t.Templates["not-found"].Execute(w, nil)
}

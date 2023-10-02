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
	Renderer.Templates["verification-mail"] = pages.VerificationMail()
	Renderer.Templates["verified"] = pages.VerifiedPage()

	Renderer.Templates["components/error"] = components.Error()
	Renderer.Templates["components/unverified-error"] = components.UnverifiedError()
	Renderer.Templates["components/urls-table"] = components.UrlsTable()
	Renderer.Templates["components/icons-row"] = components.IconsRow()
	Renderer.Templates["components/username-dialog"] = components.UsernameDialog()
	Renderer.Templates["components/password-dialog"] = components.PasswordDialog()
	Renderer.Templates["components/delete-dialog"] = components.DeleteDialog()
	Renderer.Templates["components/origin-dialog"] = components.OriginDialog()
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

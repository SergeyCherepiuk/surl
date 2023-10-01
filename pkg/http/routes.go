package http

import (
	"net/http"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/SergeyCherepiuk/surl/pkg/http/handlers"
	"github.com/SergeyCherepiuk/surl/pkg/http/middleware"
	"github.com/SergeyCherepiuk/surl/pkg/http/template"
	"github.com/SergeyCherepiuk/surl/public/views/components"
	"github.com/SergeyCherepiuk/surl/public/views/pages"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

type Router struct {
	SessionChecker domain.SessionChecker
	AccountGetter  domain.AccountGetter

	SessionCreator domain.SessionCreator
	AccountCreator domain.AccountCreator

	SessionUpdater domain.SessionUpdater
	AccountUpdater domain.AccountUpdater

	SessionDeleter domain.SessionDeleter
	AccountDeleter domain.AccountDeleter

	VerificationChecker domain.VerificationChecker
	VerificationGetter  domain.VerificationGetter
	Verificator         domain.Verificator
	VerificationDeleter domain.VerificationDeleter

	OriginGetter domain.OriginGetter
	UrlGetter    domain.UrlGetter
	UrlCreator   domain.UrlCreator
	UrlUpdater   domain.UrlUpdater
	UrlDeleter   domain.UrlDeleter
}

func (r Router) Build() *echo.Echo {
	e := echo.New()
	e.Use(echomiddleware.Logger())
	e.Static("/static", "public/static")
	e.Static("/assets", "public/assets")
	e.Renderer = template.Renderer

	// Middleware
	authMiddleware := middleware.AuthMiddleware{
		SessionChecker: r.SessionChecker,
	}
	verificationMiddleware := middleware.VerificationMiddleware{
		VerificationChecker: r.VerificationChecker,
	}
	urlMiddleware := middleware.UrlMiddleware{}

	// Handlers
	authHandler := handlers.AuthHandler{
		AccountGetter:  r.AccountGetter,
		SessionCreator: r.SessionCreator,
		AccountCreator: r.AccountCreator,
		SessionDeleter: r.SessionDeleter,
	}
	accountHandler := handlers.AccountHandler{
		AccountGetter:  r.AccountGetter,
		SessionUpdater: r.SessionUpdater,
		AccountUpdater: r.AccountUpdater,
		AccountDeleter: r.AccountDeleter,
	}
	verificationHandler := handlers.VerificationHandler{
		VerificationGetter:  r.VerificationGetter,
		Verificator:         r.Verificator,
		VerificationDeleter: r.VerificationDeleter,
	}
	urlHandler := handlers.UrlHandler{
		OriginGetter: r.OriginGetter,
		UrlGetter:    r.UrlGetter,
		UrlCreator:   r.UrlCreator,
		UrlUpdater:   r.UrlUpdater,
		UrlDeleter:   r.UrlDeleter,
	}

	// API routes
	api := e.Group("/api")

	auth := api.Group("/auth")
	auth.Use(authMiddleware.IsNotAuthenticated(func(c echo.Context) error {
		return c.NoContent(http.StatusUnauthorized)
	}))
	auth.POST("/login", authHandler.Login, verificationMiddleware.IsVerified)
	auth.POST("/signup", authHandler.SingUp)
	api.POST("/auth/signout", authHandler.SignOut, authMiddleware.IsAuthenticated(func(c echo.Context) error {
		return c.NoContent(http.StatusUnauthorized)
	}))

	account := api.Group("/account")
	account.Use(authMiddleware.IsAuthenticated(func(c echo.Context) error {
		return c.NoContent(http.StatusUnauthorized)
	}))
	account.PUT("/username", accountHandler.UpdateUsername)
	account.PUT("/password", accountHandler.UpdatePassword)
	account.DELETE("", accountHandler.Delete)

	verify := api.Group("/verify")
	verify.GET("/:username/:id", verificationHandler.Verify)

	accountViews := account.Group("/views")
	accountViews.GET("/icons-row", accountHandler.GetIconsRow)
	accountViews.GET("/username-dialog", accountHandler.GetUsernameDialog)
	accountViews.GET("/password-dialog", accountHandler.GetPasswordDialog)
	accountViews.GET("/delete-dialog", accountHandler.GetDeleteDialog)

	urls := api.Group("/urls")
	urls.Use(authMiddleware.IsAuthenticated(func(c echo.Context) error {
		return c.NoContent(http.StatusUnauthorized)
	}))
	urls.GET("", urlHandler.GetAll)
	urls.GET("/:username/:hash/listen", urlHandler.Listen)
	urls.PUT("/:username/:hash", urlHandler.Update)
	urls.POST("", urlHandler.Create)
	urls.DELETE("/:username/:hash", urlHandler.Delete, urlMiddleware.IsOwner)

	urlsViews := urls.Group("/views")
	urlsViews.GET("/origin-dialog", urlHandler.GetOriginDialog)

	// Web pages routes
	authWeb := e.Group("")
	authWeb.Use(authMiddleware.IsNotAuthenticated(func(c echo.Context) error {
		return c.Redirect(http.StatusSeeOther, "/")
	}))
	authWeb.GET("/login", func(c echo.Context) error {
		data := pages.LoginPageData{
			UsernameInputData: components.InputData{
				Type: "text", Name: "username", Placeholder: "Username", Value: "",
			},
			PasswordInputData: components.InputData{
				Type: "password", Name: "password", Placeholder: "Password", Value: "",
			},
			ButtonData: components.ButtonData{
				Type: "submit", Text: "Log in", IsPrimary: true,
			},
		}
		return c.Render(http.StatusOK, "login", data)
	})
	authWeb.GET("/signup", func(c echo.Context) error {
		data := pages.SignUpPageData{
			UsernameInputData: components.InputData{
				Type: "text", Name: "username", Placeholder: "Username", Value: "",
			},
			PasswordInputData: components.InputData{
				Type: "password", Name: "password", Placeholder: "Password", Value: "",
			},
			ButtonData: components.ButtonData{
				Type: "submit", Text: "Sing up", IsPrimary: true,
			},
		}
		return c.Render(http.StatusOK, "signup", data)
	})

	protectedWeb := e.Group("")
	protectedWeb.Use(authMiddleware.IsAuthenticated(func(c echo.Context) error {
		return c.Redirect(http.StatusSeeOther, "/login")
	}))
	protectedWeb.GET("/", func(c echo.Context) error {
		data := pages.HomePageData{
			Username: c.Get("username").(string),
		}
		return c.Render(http.StatusOK, "home", data)
	})

	e.GET("/:username/:hash", urlHandler.GetOrigin)

	e.GET("/*", func(c echo.Context) error {
		return c.Render(http.StatusOK, "not-found", nil)
	})

	return e
}

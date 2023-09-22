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

	AccountCreator domain.AccountCreator
	SessionCreator domain.SessionCreator

	AccountUpdater domain.AccountUpdater

	AccountDeleter domain.AccountDeleter

	UrlService domain.UrlService
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
	urlMiddleware := middleware.UrlMiddleware{}

	// Handlers
	userHandler := handlers.UserHandler{
		AccountGetter:  r.AccountGetter,
		AccountCreator: r.AccountCreator,
		SessionCreator: r.SessionCreator,
	}
	accountHandler := handlers.AccountHandler{
		AccountUpdater: r.AccountUpdater,
		AccountDeleter: r.AccountDeleter,
	}
	urlHandler := handlers.UrlHandler{
		UrlService: r.UrlService,
	}

	// API routes
	api := e.Group("/api")

	auth := api.Group("/auth")
	auth.Use(authMiddleware.IsNotAuthenticated(func(c echo.Context) error {
		return c.NoContent(http.StatusUnauthorized)
	}))
	auth.POST("/login", userHandler.Login)
	auth.POST("/signup", userHandler.SingUp)

	account := api.Group("/account")
	account.Use(authMiddleware.IsAuthenticated(func(c echo.Context) error {
		return c.NoContent(http.StatusUnauthorized)
	}))
	account.PUT("/username", accountHandler.UpdateUsername)
	account.DELETE("", accountHandler.Delete)

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
	urls.POST("", urlHandler.Create)
	urls.DELETE("/:username/:hash", urlHandler.Delete, urlMiddleware.IsOwner)

	// Web pages routes
	authWeb := e.Group("")
	authWeb.Use(authMiddleware.IsNotAuthenticated(func(c echo.Context) error {
		return c.Redirect(http.StatusSeeOther, "/")
	}))
	authWeb.GET("/login", func(c echo.Context) error {
		data := pages.LoginPageData{
			UsernameInputData: components.InputData{
				Type: "text", Name: "username", Placeholder: "Username",
			},
			PasswordInputData: components.InputData{
				Type: "password", Name: "password", Placeholder: "Password",
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
				Type: "text", Name: "username", Placeholder: "Username",
			},
			PasswordInputData: components.InputData{
				Type: "password", Name: "password", Placeholder: "Password",
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
			UrlInputData: components.InputWithButtonData{
				Type: "text", Name: "origin", Placeholder: "Your link", Text: "Shorten",
			},
		}
		return c.Render(http.StatusOK, "home", data)
	})

	e.GET("/:username/:hash", urlHandler.GetOrigin)

	e.GET("/*", func(c echo.Context) error {
		return c.Render(http.StatusOK, "not-found", nil)
	})

	return e
}

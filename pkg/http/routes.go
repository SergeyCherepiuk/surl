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

	VerificationSender  domain.VerificationSender
	VerificationChecker domain.VerificationChecker
	VerificationGetter  domain.VerificationGetter
	VerificationCreator domain.VerificationCreator
	Verificator         domain.Verificator
	VerificationDeleter domain.VerificationDeleter

	PasswordResetSender  domain.PasswordResetSender
	PasswordResetGetter  domain.PasswordResetGetter
	PasswordResetCreator domain.PasswordResetCreator
	PasswordResetDeleter domain.PasswordResetDeleter

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
	urlMiddleware := middleware.UrlMiddleware{}

	// Handlers
	authHandler := handlers.AuthHandler{
		AccountGetter:       r.AccountGetter,
		SessionCreator:      r.SessionCreator,
		AccountCreator:      r.AccountCreator,
		SessionDeleter:      r.SessionDeleter,
		VerificationSender:  r.VerificationSender,
		VerificationChecker: r.VerificationChecker,
		VerificationCreator: r.VerificationCreator,
	}
	accountHandler := handlers.AccountHandler{
		AccountGetter:  r.AccountGetter,
		SessionUpdater: r.SessionUpdater,
		AccountUpdater: r.AccountUpdater,
		SessionDeleter: r.SessionDeleter,
		AccountDeleter: r.AccountDeleter,
	}
	verificationHandler := handlers.VerificationHandler{
		AccountGetter:       r.AccountGetter,
		VerificationSender:  r.VerificationSender,
		VerificationGetter:  r.VerificationGetter,
		VerificationCreator: r.VerificationCreator,
		Verificator:         r.Verificator,
		VerificationDeleter: r.VerificationDeleter,
	}
	passwordResetHandler := handlers.PasswordResetHandler{
		AccountGetter:        r.AccountGetter,
		AccountUpdater:       r.AccountUpdater,
		PasswordResetSender:  r.PasswordResetSender,
		PasswordResetGetter:  r.PasswordResetGetter,
		PasswordResetCreator: r.PasswordResetCreator,
		PasswordResetDeleter: r.PasswordResetDeleter,
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
	auth.POST("/login", authHandler.Login)
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

	accountViews := account.Group("/views")
	accountViews.GET("/icons-row", accountHandler.GetIconsRow)
	accountViews.GET("/username-dialog", accountHandler.GetUsernameDialog)
	accountViews.GET("/password-dialog", accountHandler.GetPasswordDialog)
	accountViews.GET("/delete-dialog", accountHandler.GetDeleteDialog)

	verification := api.Group("/verification")
	verification.GET("/:username/:id", verificationHandler.Verify)
	verification.POST("/send/:username", verificationHandler.Send)

	passwordReset := api.Group("/password-reset")
	passwordReset.GET("/:username/:id", passwordResetHandler.GetPasswordResetPage)
	passwordReset.POST("/send", passwordResetHandler.Send)
	passwordReset.POST("/reset/:username", passwordResetHandler.Reset)

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
			EmailInputData: components.InputData{
				Type: "email", Name: "email", Placeholder: "Email", Value: "",
			},
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
	authWeb.GET("/forgot-password", func(c echo.Context) error {
		data := pages.ForgotPasswordPageData{
			UsernameInputData: components.InputData{
				Type: "text", Name: "username", Placeholder: "Username", Value: "",
			},
			ButtonData: components.ButtonData{
				Type: "submit", Text: "Send mail", IsPrimary: true,
			},
		}
		return c.Render(http.StatusOK, "forgot-password", data)
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

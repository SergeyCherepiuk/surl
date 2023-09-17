package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/SergeyCherepiuk/surl/pkg/database/postgres"
	"github.com/SergeyCherepiuk/surl/pkg/database/redis"
	"github.com/SergeyCherepiuk/surl/pkg/http/handlers"
	"github.com/SergeyCherepiuk/surl/pkg/http/middleware"
	"github.com/SergeyCherepiuk/surl/pkg/http/template"
	"github.com/SergeyCherepiuk/surl/public/views/components"
	"github.com/SergeyCherepiuk/surl/public/views/pages"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Panic(err)
	}
	postgres.MustConnect()
	redis.MustConnect()
}

func main() {
	e := echo.New()
	e.Use(echomiddleware.Logger())
	e.Static("/assets", "public/assets")
	e.Renderer = template.Renderer

	// Services
	accountManagerService := postgres.NewAccountManagerService()
	sessionManagerService := redis.NewSessionManagerService()
	urlService := postgres.NewUrlService()

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(sessionManagerService)

	// Handlers
	userHandler := handlers.UserHandler{
		AccountManagerService: accountManagerService,
		SessionManagerService: sessionManagerService,
	}
	urlHandler := handlers.UrlHandler{
		UrlService: urlService,
	}

	// API routes
	api := e.Group("/api")

	auth := api.Group("/auth")
	auth.POST("/login", userHandler.Login)
	auth.POST("/signup", userHandler.SingUp)

	urls := api.Group("/urls")
	urls.Use(authMiddleware.CheckIfAuthenticated) // TODO: Check middleware (redirection might work incorrectly)
	urls.GET("", urlHandler.GetAll)
	urls.POST("", urlHandler.Create)

	// Web pages routes
	authWeb := e.Group("")
	authWeb.GET("/login", func(c echo.Context) error {
		data := pages.LoginPageData{
			UsernameInputData: components.InputComponentData{
				Type: "text", Name: "username", Placeholder: "Username",
			},
			PasswordInputData: components.InputComponentData{
				Type: "password", Name: "password", Placeholder: "Password",
			},
			ButtonData: components.ButtonComponentData{
				Type: "submit", Text: "Log in",
			},
		}
		return c.Render(http.StatusOK, "login", data)
	})
	authWeb.GET("/signup", func(c echo.Context) error {
		data := pages.SignUpPageData{
			UsernameInputData: components.InputComponentData{
				Type: "text", Name: "username", Placeholder: "Username",
			},
			PasswordInputData: components.InputComponentData{
				Type: "password", Name: "password", Placeholder: "Password",
			},
			ButtonData: components.ButtonComponentData{
				Type: "submit", Text: "Sing up",
			},
		}
		return c.Render(http.StatusOK, "signup", data)
	})

	protectedWeb := e.Group("")
	protectedWeb.Use(authMiddleware.CheckIfAuthenticated)
	protectedWeb.GET("/", func(c echo.Context) error {
		data := pages.HomePageData{
			Username: c.Get("username").(string),
			UrlInputData: components.InputWithButtonComponentData{
				Type: "text", Name: "origin", Placeholder: "Your link", Text: "Shorten",
			},
		}
		return c.Render(http.StatusOK, "home", data)
	})

	e.GET("/:username/:hash", urlHandler.GetOrigin)

	e.Start(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
}

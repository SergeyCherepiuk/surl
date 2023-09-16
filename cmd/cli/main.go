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
	auth.Use(authMiddleware.CheckIfNotAuthenticated)
	auth.POST("/login", userHandler.Login)
	auth.POST("/signup", userHandler.SingUp)

	urls := api.Group("/urls")
	urls.Use(authMiddleware.CheckIfAuthenticated)
	urls.GET("", urlHandler.GetAll)
	urls.POST("", urlHandler.Create)

	// Web pages routes
	authWeb := e.Group("")
	authWeb.Use(authMiddleware.CheckIfNotAuthenticated)
	authWeb.GET("/login", func(c echo.Context) error {
		return c.Render(http.StatusOK, "login", nil)
	})
	authWeb.GET("/signup", func(c echo.Context) error {
		return c.Render(http.StatusOK, "signup", nil)
	})

	protectedWeb := e.Group("")
	protectedWeb.Use(authMiddleware.CheckIfAuthenticated)
	protectedWeb.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "home", c.Get("username").(string))
	})

	e.Start(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
}

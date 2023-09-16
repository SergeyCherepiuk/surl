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

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(sessionManagerService)

	// Handlers
	userHandler := handlers.UserHandler{
		AccountManagerService: accountManagerService,
		SessionManagerService: sessionManagerService,
	}

	// API routes
	api := e.Group("/api")

	api.POST("/auth/login", userHandler.Login)
	api.POST("/auth/signup", userHandler.SingUp)

	// Web pages routes
	auth := e.Group("")
	auth.Use(authMiddleware.CheckIfNotAuthenticated)
	auth.GET("/login", func(c echo.Context) error {
		return c.Render(http.StatusOK, "login", nil)
	})
	auth.GET("/signup", func(c echo.Context) error {
		return c.Render(http.StatusOK, "signup", nil)
	})

	protected := e.Group("")
	protected.Use(authMiddleware.CheckIfAuthenticated)
	protected.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "home", nil)
	})

	e.Start(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
}

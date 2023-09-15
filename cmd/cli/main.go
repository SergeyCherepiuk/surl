package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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

	accountManagerService := postgres.NewAccountManagerService()
	sessionManagerService := redis.NewSessionManagerService()

	userHandler := handlers.UserHandler{
		AccountManagerService: accountManagerService,
		SessionManagerService: sessionManagerService,
	}

	api := e.Group("/api")
	api.POST("/auth/signup", userHandler.SingUp)

	e.GET("/:page", func(c echo.Context) error {
		page := strings.TrimPrefix(c.Param("page"), "components/")
		return c.Render(http.StatusOK, page, nil)
	})

	authProtected := e.Group("")
	authProtected.Use(middleware.NewAuthMiddleware(sessionManagerService).Check)
	authProtected.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "home", nil)
	})

	e.Start(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
}

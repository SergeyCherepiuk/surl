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
	"github.com/SergeyCherepiuk/surl/pkg/http/template"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	e.Use(middleware.Logger())
	e.Renderer = template.Renderer

	userHandler := handlers.UserHandler{
		AccountManagerService: postgres.NewAccountManagerService(),
		SessionManagerService: redis.NewSessionManagerService(),
	}

	api := e.Group("/api")
	api.POST("/auth/signup", userHandler.SingUp)

	e.GET("/:page", func(c echo.Context) error {
		page := strings.TrimPrefix(c.Param("page"), "components/")
		return c.Render(http.StatusOK, page, nil)
	})

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "home", nil)
	})

	e.Start(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
}

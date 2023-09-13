package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/SergeyCherepiuk/surl/pkg/http/template"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Panic(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Renderer = template.Renderer

	e.GET("/:page", func(c echo.Context) error {
		page := strings.TrimPrefix(c.Param("page"), "components/")
		return c.Render(http.StatusOK, page, nil)
	})

	e.Start(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
}

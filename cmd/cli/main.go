package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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

	e.GET("/:name", func(c echo.Context) error {
		data := struct{ Name string }{Name: c.Param("name")}
		return c.Render(http.StatusOK, "greeting.html", data)
	})

	e.GET("/*", func(c echo.Context) error {
		return c.Render(http.StatusNotFound, "404.html", nil)
	})

	e.Start(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
}

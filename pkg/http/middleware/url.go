package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type UrlMiddleware struct{}

func (um UrlMiddleware) IsOwner(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Get("username").(string) != c.Param("username") {
			return c.NoContent(http.StatusUnauthorized)
		}

		return next(c)
	}
}

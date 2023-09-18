package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type urlMiddleware struct{}

func NewUrlMiddleware() *urlMiddleware {
	return &urlMiddleware{}
}

func (m urlMiddleware) IsOwner(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Get("username").(string) != c.Param("username") {
			return c.NoContent(http.StatusUnauthorized)
		}

		return next(c)
	}
}

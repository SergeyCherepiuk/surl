package middleware

import (
	"context"
	"net/http"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/labstack/echo/v4"
)

type VerificationMiddleware struct {
	VerificationChecker domain.VerificationChecker
}

func (m VerificationMiddleware) IsVerified(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.FormValue("username")

		if err := m.VerificationChecker.Check(context.Background(), username); err != nil {
			return c.NoContent(http.StatusUnauthorized)
		}

		return next(c)
	}
}

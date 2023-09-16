package middleware

import (
	"context"
	"net/http"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type authMiddleware struct {
	sessionManagerService domain.SessionManagerService
}

func NewAuthMiddleware(sessionManagerService domain.SessionManagerService) *authMiddleware {
	return &authMiddleware{sessionManagerService: sessionManagerService}
}

func (m authMiddleware) CheckIfAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_id")
		if err != nil {
			return c.Redirect(http.StatusMovedPermanently, "/login")
		}

		id, err := uuid.Parse(cookie.Value)
		if err != nil {
			return c.Redirect(http.StatusMovedPermanently, "/login")
		}

		username, err := m.sessionManagerService.Check(context.Background(), id)
		if err != nil {
			return c.Redirect(http.StatusMovedPermanently, "/login")
		}

		c.Set("username", username)
		return next(c)
	}
}

func (m authMiddleware) CheckIfNotAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_id")
		if err != nil {
			return next(c)
		}

		id, err := uuid.Parse(cookie.Value)
		if err != nil {
			return next(c)
		}

		username, err := m.sessionManagerService.Check(context.Background(), id)
		if err != nil {
			return next(c)
		}

		c.Set("username", username)
		return c.Redirect(http.StatusMovedPermanently, "/")
	}
}

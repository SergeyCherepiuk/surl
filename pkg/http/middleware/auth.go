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

func (m authMiddleware) Check(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_id")
		if err != nil {
			return c.Redirect(http.StatusMovedPermanently, "/login")
		}

		id, err := uuid.Parse(cookie.Value)
		if err != nil {
			return c.Redirect(http.StatusMovedPermanently, "/login")
		}

		if err := m.sessionManagerService.Check(context.Background(), id); err != nil {
			return c.Redirect(http.StatusMovedPermanently, "/login")
		}

		return next(c)
	}
}

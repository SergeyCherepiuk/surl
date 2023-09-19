package middleware

import (
	"context"

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

func (m authMiddleware) IsAuthenticated(onError echo.HandlerFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("session_id")
			if err != nil {
				return onError(c)
			}

			id, err := uuid.Parse(cookie.Value)
			if err != nil {
				return onError(c)
			}

			username, err := m.sessionManagerService.Check(context.Background(), id)
			if err != nil {
				return onError(c)
			}

			c.Set("username", username)
			return next(c)
		}
	}
}

func (m authMiddleware) IsNotAuthenticated(onError echo.HandlerFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
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
			return onError(c)
		}
	}
}

package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/SergeyCherepiuk/surl/pkg/http/validation"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	SessionManagerService domain.SessionManagerService
	AccountManagerService domain.AccountManagerService
}

func (h UserHandler) SingUp(c echo.Context) error {
	user := domain.User{
		Username: c.FormValue("username"),
		Password: c.FormValue("password"),
	}

	if err := validation.ValidateUserSignUp(user); err != nil {
		return c.Render(http.StatusOK, "component/error", err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Username), 10)
	if err != nil {
		return c.Render(http.StatusOK, "component/error", "Failed to encrypt your password. Please try again.")
	}
	user.Password = string(hashedPassword)

	if err := h.AccountManagerService.Create(context.Background(), user); err != nil {
		return c.Render(http.StatusOK, "component/error", "Failed to save your account to the database. Please try again.")
	}

	ttl := 7 * 24 * time.Hour
	id, err := h.SessionManagerService.Create(context.Background(), user.Username, ttl)
	if err != nil {
		c.Response().Header().Set("HX-Redirect", "/login")
		return c.NoContent(http.StatusMovedPermanently)
	}

	c.SetCookie(&http.Cookie{
		Name:     "session_id",
		Value:    id.String(),
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(ttl),
	})
	c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(http.StatusMovedPermanently)
}

package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/SergeyCherepiuk/surl/pkg/http/validation"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	SessionManagerService domain.SessionManagerService
	AccountManagerService domain.AccountManagerService
}

func (h UserHandler) Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if err := validation.ValidateAuthentication(
		domain.User{
			Username: c.FormValue("username"),
			Password: c.FormValue("password"),
		},
	); err != nil {
		return c.Render(http.StatusOK, "component/error", err.Error())
	}

	user, err := h.AccountManagerService.Get(context.Background(), username)
	if err != nil {
		return c.Render(http.StatusOK, "component/error", "No user with this username was found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.Logger().Error(err)
		return c.Render(http.StatusOK, "component/error", "Wrong password")
	}

	ttl := 7 * 24 * time.Hour
	id, err := h.SessionManagerService.Create(context.Background(), user.Username, ttl)
	if err != nil {
		c.Response().Header().Set("HX-Redirect", "/login")
		return c.NoContent(http.StatusMovedPermanently)
	}

	h.setCookie(c, id, ttl)
	c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(http.StatusOK)
}

func (h UserHandler) SingUp(c echo.Context) error {
	user := domain.User{
		Username: c.FormValue("username"),
		Password: c.FormValue("password"),
	}

	if err := validation.ValidateAuthentication(user); err != nil {
		return c.Render(http.StatusOK, "component/error", err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
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

	h.setCookie(c, id, ttl)
	c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(http.StatusMovedPermanently)
}

func (h UserHandler) setCookie(c echo.Context, id uuid.UUID, ttl time.Duration) {
	c.SetCookie(&http.Cookie{
		Name:     "session_id",
		Value:    id.String(),
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(ttl),
	})
}

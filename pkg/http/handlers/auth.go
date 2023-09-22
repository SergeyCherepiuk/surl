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
	AccountGetter  domain.AccountGetter
	SessionCreator domain.SessionCreator
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
		return c.String(http.StatusOK, err.Error())
	}

	user, err := h.AccountGetter.Get(context.Background(), username)
	if err != nil {
		return c.String(http.StatusOK, "No user with this username was found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return c.String(http.StatusOK, "Wrong password")
	}

	ttl := 7 * 24 * time.Hour
	id, err := h.SessionCreator.Create(context.Background(), user, ttl)
	if err != nil {
		return c.String(http.StatusOK, "Failed to create a session")
	}

	h.setCookie(c, id, ttl)
	c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(http.StatusSeeOther)
}

func (h UserHandler) SingUp(c echo.Context) error {
	user := domain.User{
		Username: c.FormValue("username"),
		Password: c.FormValue("password"),
	}

	if err := validation.ValidateAuthentication(user); err != nil {
		return c.String(http.StatusOK, err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return c.String(http.StatusOK, "Failed to encrypt your password. Please try again")
	}
	user.Password = string(hashedPassword)

	ttl := 7 * 24 * time.Hour
	id, err := h.SessionCreator.Create(context.Background(), user, ttl)
	if err != nil {
		c.Response().Header().Set("HX-Redirect", "/login")
		return c.NoContent(http.StatusSeeOther)
	}

	h.setCookie(c, id, ttl)
	c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(http.StatusSeeOther)
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

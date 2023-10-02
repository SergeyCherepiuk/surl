package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/SergeyCherepiuk/surl/pkg/http/validation"
	"github.com/SergeyCherepiuk/surl/pkg/mail"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	AccountGetter       domain.AccountGetter
	AccountCreator      domain.AccountCreator
	SessionCreator      domain.SessionCreator
	SessionDeleter      domain.SessionDeleter
	VerificationChecker domain.VerificationChecker
	VerificationCreator domain.VerificationCreator
}

func (h AuthHandler) Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if err := validation.ValidateUsername(username); err != nil {
		return c.Render(http.StatusOK, "components/error", err.Error())
	} else if err := validation.ValidatePassword(password); err != nil {
		return c.Render(http.StatusOK, "components/error", err.Error())
	}

	user, err := h.AccountGetter.Get(context.Background(), username)
	if err != nil {
		return c.Render(http.StatusOK, "components/error", "No user with this username was found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return c.Render(http.StatusOK, "components/error", "Wrong password")
	}

	if err := h.VerificationChecker.Check(context.Background(), user.Username); err != nil {
		return c.Render(http.StatusOK, "components/error", "Account is not verified")
	}

	ttl := 7 * 24 * time.Hour
	id, err := h.SessionCreator.Create(context.Background(), user.Username, ttl)
	if err != nil {
		return c.Render(http.StatusOK, "components/error", "Failed to create a session")
	}

	h.setCookie(c, id, ttl)
	c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(http.StatusSeeOther)
}

func (h AuthHandler) SingUp(c echo.Context) error {
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")

	if err := validation.ValidateUsername(username); err != nil {
		return c.Render(http.StatusOK, "components/error", err.Error())
	} else if err := validation.ValidatePassword(password); err != nil {
		return c.Render(http.StatusOK, "components/error", err.Error())
	}

	user := domain.User{
		Username: username,
		Email:    email,
		Password: password,
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return c.Render(http.StatusOK, "components/error", "Failed to encrypt your password. Please try again")
	}
	user.Password = string(hashedPassword)

	if err := h.AccountCreator.Create(context.Background(), user); err != nil {
		return c.Render(http.StatusOK, "components/error", "Failed save your data to the database")
	}

	verificationRequestId := uuid.NewString()
	if err := h.VerificationCreator.Create(context.Background(), username, verificationRequestId); err != nil {
		return c.Render(http.StatusOK, "components/error", "Failed to send verification email try to login and request a new one")
	}

	go func() {
		// TODO: Send simple html formatted page instead
		verificationLink := fmt.Sprintf("http://localhost:3000/api/verify/%s/%s", username, verificationRequestId)
		mail.Sender.Send(email, "Account verification", verificationLink)
	}()

	c.Response().Header().Set("HX-Redirect", "/login")
	return c.NoContent(http.StatusSeeOther)
}

func (h AuthHandler) SignOut(c echo.Context) error {
	username := c.Get("username").(string)

	if err := h.SessionDeleter.Delete(context.Background(), username); err != nil {
		return c.Render(http.StatusOK, "components/error", "Failed to invalidate the session")
	}

	c.SetCookie(&http.Cookie{
		Name:    "session_id",
		Value:   "",
		Path:    "/",
		Expires: time.Now(),
	})
	c.Response().Header().Set("HX-Redirect", "/login")
	return c.NoContent(http.StatusOK)
}

func (h AuthHandler) setCookie(c echo.Context, id uuid.UUID, ttl time.Duration) {
	c.SetCookie(&http.Cookie{
		Name:     "session_id",
		Value:    id.String(),
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(ttl),
	})
}

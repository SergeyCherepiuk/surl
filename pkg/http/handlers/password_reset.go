package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/SergeyCherepiuk/surl/pkg/http/validation"
	"github.com/SergeyCherepiuk/surl/public/views/components"
	"github.com/SergeyCherepiuk/surl/public/views/pages"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type PasswordResetHandler struct {
	AccountGetter        domain.AccountGetter
	AccountUpdater       domain.AccountUpdater
	PasswordResetSender  domain.PasswordResetSender
	PasswordResetGetter  domain.PasswordResetGetter
	PasswordResetCreator domain.PasswordResetCreator
	PasswordResetDeleter domain.PasswordResetDeleter
}

func (h PasswordResetHandler) GetPasswordResetPage(c echo.Context) error {
	username := c.Param("username")
	id := c.Param("id")

	passwordResetRequest, err := h.PasswordResetGetter.Get(context.Background(), username, id)
	if err != nil {
		return c.Render(http.StatusOK, "not-found", nil)
	}

	if passwordResetRequest.ExpiresAt.Before(time.Now().In(time.UTC)) {
		return c.Render(http.StatusOK, "not-found", nil)
	}

	data := pages.PasswordResetPageData{
		Username: username,
		NewPasswordInputData: components.InputData{
			Type: "password", Name: "new-password", Placeholder: "New password", Value: "",
		},
		NewPasswordRepeatInputData: components.InputData{
			Type: "password", Name: "new-password-repeat", Placeholder: "New password (repeat)", Value: "",
		},
		ButtonData: components.ButtonData{
			Type: "submit", Text: "Reset password", IsPrimary: true,
		},
	}
	return c.Render(http.StatusOK, "password-reset", data)
}

func (h PasswordResetHandler) Reset(c echo.Context) error {
	username := c.Param("username")

	body := make(map[string]any)
	if err := c.Bind(&body); err != nil {
		return c.Render(http.StatusOK, "components/error", "Invalid request body")
	}

	newPassword, newPasswordOk := body["new-password"].(string)
	newPasswordRepeat, newPasswordRepeatOk := body["new-password-repeat"].(string)
	if !newPasswordOk || !newPasswordRepeatOk {
		return c.Render(http.StatusOK, "components/error", "Failed to get password from the request")
	}

	if err := validation.ValidatePassword(newPassword); err != nil {
		return c.Render(http.StatusOK, "components/error", err.Error())
	}

	if newPassword != newPasswordRepeat {
		return c.Render(http.StatusOK, "components/error", "New passwords aren't the same")
	}

	newPasswordHashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	if err != nil {
		return c.Render(http.StatusOK, "components/error", "Failed to hash new password")
	}

	if err := h.AccountUpdater.UpdatePassword(context.Background(), username, string(newPasswordHashed)); err != nil {
		return c.Render(http.StatusOK, "components/error", "Failed to update the password")
	}

	h.PasswordResetDeleter.DeleteAll(context.Background(), username)

	c.Response().Header().Set("HX-Redirect", "/login")
	return c.NoContent(http.StatusOK)
}

func (h PasswordResetHandler) Send(c echo.Context) error {
	username := c.FormValue("username")

	user, err := h.AccountGetter.Get(context.Background(), username)
	if err != nil {
		return c.String(http.StatusBadRequest, "No user with this username was found")
	}

	passwordResetRequest := domain.PasswordResetRequest{
		ID:        uuid.NewString(),
		Username:  user.Username,
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}

	if err := h.PasswordResetCreator.Create(context.Background(), passwordResetRequest); err != nil {
		return c.String(http.StatusBadRequest, "Failed to send a password reset mail")
	}

	go func() {
		h.PasswordResetSender.Send(context.Background(), user.Email, passwordResetRequest)
	}()

	return c.Render(http.StatusOK, "components/info", "Password reset request has been sent")
}

package handlers

import (
	"context"
	"net/http"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type VerificationHandler struct {
	AccountGetter       domain.AccountGetter
	VerificationSender  domain.VerificationSender
	VerificationGetter  domain.VerificationGetter
	VerificationCreator domain.VerificationCreator
	Verificator         domain.Verificator
	VerificationDeleter domain.VerificationDeleter
}

func (h VerificationHandler) Verify(c echo.Context) error {
	username := c.Param("username")
	id := c.Param("id")

	// TODO: Use "verificationRequest" to check if it's still valid (after adding expiration time to it)
	if _, err := h.VerificationGetter.Get(context.Background(), username, id); err != nil {
		return c.Render(http.StatusOK, "not-found", nil)
	}

	if err := h.Verificator.Verify(context.Background(), username); err != nil {
		return c.Render(http.StatusOK, "not-found", nil)
	}

	h.VerificationDeleter.DeleteAll(context.Background(), username) // NOTE: Error is ignored
	return c.Render(http.StatusOK, "verified", nil)
}

func (h VerificationHandler) Send(c echo.Context) error {
	username := c.Param("username")

	user, err := h.AccountGetter.Get(context.Background(), username)
	if err != nil {
		return c.String(http.StatusBadRequest, "No user with this username was found")
	}

	if user.IsVerified {
		return c.String(http.StatusBadRequest, "User is already verified")
	}

	verificationRequestId := uuid.NewString()
	if err := h.VerificationCreator.Create(context.Background(), username, verificationRequestId); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to send verification email try to login and request a new one")
	}

	go func() {
		h.VerificationSender.Send(user.Email, username, verificationRequestId)
	}()

	return c.NoContent(http.StatusOK)
}

package handlers

import (
	"context"
	"net/http"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type VerificationHandler struct {
	VerificationGetter  domain.VerificationGetter
	Verificator         domain.Verificator
	VerificationDeleter domain.VerificationDeleter
}

func (vh VerificationHandler) Verify(c echo.Context) error {
	username := c.Param("username")
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.Render(http.StatusOK, "not-found", nil)
	}

	// TODO: Use "verificationRequest" to check if it's still valid (after adding expiration time to it)
	if _, err := vh.VerificationGetter.Get(context.Background(), username, id); err != nil {
		return c.Render(http.StatusOK, "not-found", nil)
	}

	if err := vh.Verificator.Verify(context.Background(), username); err != nil {
		return c.Render(http.StatusOK, "not-found", nil)
	}

	vh.VerificationDeleter.DeleteAll(context.Background(), username) // NOTE: Error is ignored
	return c.Render(http.StatusOK, "verified", nil)
}

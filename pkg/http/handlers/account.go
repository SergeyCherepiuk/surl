package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/SergeyCherepiuk/surl/pkg/http/validation"
	"github.com/SergeyCherepiuk/surl/public/views/components"
	"github.com/labstack/echo/v4"
)

type AccountHandler struct {
	AccountUpdater domain.AccountUpdater
	AccountDeleter domain.AccountDeleter
}

func (h AccountHandler) GetIconsRow(c echo.Context) error {
	data := components.IconsRowComponentData{
		ChangeUsernameIconButtonData: components.IconButtonComponentData{
			Type: "button", Icon: "assets/images/ic-edit.svg", Alt: "Change username",
		},
		ChangePasswordIconButtonData: components.IconButtonComponentData{
			Type: "button", Icon: "assets/images/ic-security.svg", Alt: "Change password",
		},
		DeleteAccountIconButtonData: components.IconButtonComponentData{
			Type: "button", Icon: "assets/images/ic-delete.svg", Alt: "Delete account",
		},
	}
	return c.Render(http.StatusOK, "components/icons-row", data)
}

func (h AccountHandler) GetUsernameDialog(c echo.Context) error {
	data := components.UsernameDialogComponentData{
		InputData: components.InputComponentData{
			Type: "text", Name: "new-username", Placeholder: "New username",
		},
		ConfirmIconButtonData: components.IconButtonComponentData{
			Type: "submit", Icon: "assets/images/ic-confirm.svg", Alt: "Confirm",
		},
		DeclineIconButtonData: components.IconButtonComponentData{
			Type: "button", Icon: "assets/images/ic-decline.svg", Alt: "Decline",
		},
	}
	return c.Render(http.StatusOK, "components/username-dialog", data)
}

func (h AccountHandler) GetDeleteDialog(c echo.Context) error {
	data := components.DeleteDialogComponentData{
		Text: "Are you sure you want to delete the account?",
		ConfirmIconButtonData: components.IconButtonComponentData{
			Type: "submit", Icon: "assets/images/ic-confirm.svg", Alt: "Confirm",
		},
		DeclineIconButtonData: components.IconButtonComponentData{
			Type: "button", Icon: "assets/images/ic-decline.svg", Alt: "Decline",
		},
	}
	return c.Render(http.StatusOK, "components/delete-dialog", data)
}

func (h AccountHandler) UpdateUsername(c echo.Context) error {
	username := c.Get("username").(string)
	newUsername := c.FormValue("new-username")

	if err := validation.ValidateUsernameChange(newUsername); err != nil {
		return c.String(http.StatusOK, err.Error())
	}

	if err := h.AccountUpdater.UpdateUsername(context.Background(), username, newUsername); err != nil {
		return c.String(http.StatusOK, "Failed to update username in the database")
	}

	c.Response().Header().Set("HX-Refresh", "true")
	return h.GetIconsRow(c)
}

func (h AccountHandler) Delete(c echo.Context) error {
	username := c.Get("username").(string)

	if err := h.AccountDeleter.Delete(context.Background(), username); err != nil {
		return c.String(http.StatusOK, "Failed to delete the account")
	}

	c.SetCookie(&http.Cookie{
		Name:    "session_id",
		Value:   "",
		Path:    "/",
		Expires: time.Now(),
	})
	c.Response().Header().Set("HX-Redirect", "/signup")
	return c.NoContent(http.StatusSeeOther)
}

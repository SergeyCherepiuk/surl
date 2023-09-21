package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/SergeyCherepiuk/surl/public/views/components"
	"github.com/labstack/echo/v4"
)

type AccountHandler struct {
	AccountManagerService domain.AccountManagerService
	SessionManagerService domain.SessionManagerService
}

func (h AccountHandler) GetDeleteDialog(c echo.Context) error {
	data := components.DialogComponentData{
		Text: "Are you sure you want to delete the account?",
		ConfirmIconButtonData: components.IconButtonComponentData{
			Icon: "assets/images/ic-confirm.svg", Alt: "Confirm",
		},
		DeclineIconButtonData: components.IconButtonComponentData{
			Icon: "assets/images/ic-decline.svg", Alt: "Decline",
		},
	}
	return c.Render(http.StatusOK, "components/dialog", data)
}

func (h AccountHandler) GetIconsRow(c echo.Context) error {
	data := components.IconsRowComponentData{
		ChangeUsernameIconButtonData: components.IconButtonComponentData{
			Icon: "assets/images/ic-edit.svg", Alt: "Change username",
		},
		ChangePasswordIconButtonData: components.IconButtonComponentData{
			Icon: "assets/images/ic-security.svg", Alt: "Change password",
		},
		DeleteAccountIconButtonData: components.IconButtonComponentData{
			Icon: "assets/images/ic-delete.svg", Alt: "Delete account",
		},
	}
	return c.Render(http.StatusOK, "components/icons-row", data)
}

func (h AccountHandler) Delete(c echo.Context) error {
	username := c.Get("username").(string)

	if err := h.AccountManagerService.Delete(context.Background(), username); err != nil {
		return c.String(http.StatusOK, "Failed to delete the account")
	}

	h.SessionManagerService.Invalidate(context.Background(), username)

	c.SetCookie(&http.Cookie{
		Name:    "session_id",
		Value:   "",
		Path:    "/",
		Expires: time.Now(),
	})
	c.Response().Header().Set("HX-Redirect", "/signup")
	return c.NoContent(http.StatusSeeOther)
}

package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/SergeyCherepiuk/surl/pkg/http/validation"
	"github.com/SergeyCherepiuk/surl/public/views/components"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type AccountHandler struct {
	AccountGetter  domain.AccountGetter
	AccountUpdater domain.AccountUpdater
	SessionUpdater domain.SessionUpdater
	AccountDeleter domain.AccountDeleter
}

func (h AccountHandler) GetIconsRow(c echo.Context) error {
	data := components.IconsRowData{
		ChangeUsernameIconButtonData: components.IconButtonData{
			Type: "button", Icon: "assets/images/ic-edit.svg", Alt: "Change username",
		},
		ChangePasswordIconButtonData: components.IconButtonData{
			Type: "button", Icon: "assets/images/ic-security.svg", Alt: "Change password",
		},
		DeleteAccountIconButtonData: components.IconButtonData{
			Type: "button", Icon: "assets/images/ic-delete.svg", Alt: "Delete account",
		},
	}
	return c.Render(http.StatusOK, "components/icons-row", data)
}

func (h AccountHandler) GetUsernameDialog(c echo.Context) error {
	data := components.UsernameDialogData{
		InputData: components.InputData{
			Type: "text", Name: "new-username", Placeholder: "New username",
		},
		ConfirmIconButtonData: components.IconButtonData{
			Type: "submit", Icon: "assets/images/ic-confirm.svg", Alt: "Confirm",
		},
		DeclineIconButtonData: components.IconButtonData{
			Type: "button", Icon: "assets/images/ic-decline.svg", Alt: "Decline",
		},
	}
	return c.Render(http.StatusOK, "components/username-dialog", data)
}

func (h AccountHandler) GetPasswordDialog(c echo.Context) error {
	data := components.PasswordDialogData{
		OldPasswordInputData: components.InputData{
			Type: "password", Name: "old-password", Placeholder: "Old password",
		},
		NewPasswordInputData: components.InputData{
			Type: "password", Name: "new-password", Placeholder: "New password",
		},
		NewPasswordRepeatInputData: components.InputData{
			Type: "password", Name: "new-password-repeat", Placeholder: "New password (repeat)",
		},
		SubmitButtonData: components.ButtonData{
			Type: "submit", Text: "Change", IsPrimary: true,
		},
		GoBackButtonData: components.ButtonData{
			Type: "button", Text: "Go back", IsPrimary: false,
		},
	}
	return c.Render(http.StatusOK, "components/password-dialog", data)
}

func (h AccountHandler) GetDeleteDialog(c echo.Context) error {
	data := components.DeleteDialogData{
		Text: "Are you sure you want to delete the account?",
		ConfirmIconButtonData: components.IconButtonData{
			Type: "submit", Icon: "assets/images/ic-confirm.svg", Alt: "Confirm",
		},
		DeclineIconButtonData: components.IconButtonData{
			Type: "button", Icon: "assets/images/ic-decline.svg", Alt: "Decline",
		},
	}
	return c.Render(http.StatusOK, "components/delete-dialog", data)
}

func (h AccountHandler) UpdateUsername(c echo.Context) error {
	username := c.Get("username").(string)
	newUsername := c.FormValue("new-username")

	if err := validation.ValidateUsername(newUsername); err != nil {
		return c.String(http.StatusOK, err.Error())
	}

	// TODO: AccountUpdater.UpdateUsername is called implicitly here
	if err := h.SessionUpdater.UpdateUsername(context.Background(), username, newUsername); err != nil {
		return c.String(http.StatusOK, "Failed to update username in the database")
	}

	c.Response().Header().Set("HX-Refresh", "true")
	return h.GetIconsRow(c)
}

func (h AccountHandler) UpdatePassword(c echo.Context) error {
	username := c.Get("username").(string)
	oldPassword := c.FormValue("old-password")
	newPassword := c.FormValue("new-password")
	newPasswordRepeat := c.FormValue("new-password-repeat")

	if err := validation.ValidatePassword(oldPassword); err != nil {
		return c.String(http.StatusOK, err.Error())
	} else if err := validation.ValidatePassword(newPassword); err != nil {
		return c.String(http.StatusOK, err.Error())
	}

	user, err := h.AccountGetter.Get(context.Background(), username)
	if err != nil {
		return c.String(http.StatusOK, "Failed to find your account in the database")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return c.String(http.StatusOK, "Wrong old password")
	}

	if newPassword != newPasswordRepeat {
		return c.String(http.StatusOK, "New passwords aren't the same")
	}

	newPasswordHashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	if err != nil {
		return c.String(http.StatusOK, "Failed to hash new password")
	}

	if err := h.AccountUpdater.UpdatePassword(context.Background(), username, string(newPasswordHashed)); err != nil {
		return c.String(http.StatusOK, "Failed to update the password")
	}

	return c.NoContent(http.StatusOK)
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

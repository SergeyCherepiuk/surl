package handlers

import (
	"context"
	"fmt"
	"hash/crc32"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/SergeyCherepiuk/surl/pkg/http/sse"
	"github.com/SergeyCherepiuk/surl/pkg/http/validation"
	"github.com/SergeyCherepiuk/surl/public/views/components"
	"github.com/labstack/echo/v4"
)

type UrlHandler struct {
	OriginGetter domain.OriginGetter
	UrlGetter    domain.UrlGetter
	UrlCreator   domain.UrlCreator
	UrlUpdater   domain.UrlUpdater
	UrlDeleter   domain.UrlDeleter
}

func (h UrlHandler) GetOriginDialog(c echo.Context) error {
	body := make(map[string]any)
	if err := c.Bind(&body); err != nil {
		return c.Render(http.StatusOK, "components/error", "Invalid request body")
	}

	username, usernameOk := body["username"].(string)
	hash, hashOk := body["hash"].(string)
	origin, originOk := body["origin"].(string)
	if !usernameOk || !hashOk || !originOk {
		return c.Render(http.StatusOK, "components/error", "Failed to get url details from the request")
	}

	data := components.OriginDialogData{
		Username: username,
		Hash:     hash,
		OriginInputData: components.InputData{
			Type: "text", Name: "new-origin", Placeholder: "New origin link", Value: origin,
		},
		SubmitButtonData: components.ButtonData{
			Type: "submit", Text: "Save", IsPrimary: true,
		},
		GoBackButtonData: components.ButtonData{
			Type: "button", Text: "Go back", IsPrimary: false,
		},
	}

	return c.Render(http.StatusOK, "components/origin-dialog", data)
}

// BUG: Erases the link if called on expired link
func (h UrlHandler) GetOrigin(c echo.Context) error {
	c.Response().Header().Set("Cache-Control", "no-cache, max-age=0")

	username := c.Param("username")
	hash := c.Param("hash")

	origin, _, err := h.OriginGetter.Get(context.Background(), username, hash)
	if err != nil {
		return c.Render(http.StatusOK, "not-found", nil)
	}

	updates := domain.UrlUpdates{
		Origin:     origin,
		Hash:       hash,
		LastUsedAt: time.Now(),
	}
	h.UrlUpdater.Update(context.Background(), username, hash, updates) // NOTE: Error is ignored
	return c.Redirect(http.StatusMovedPermanently, origin)
}

func (h UrlHandler) GetAll(c echo.Context) error {
	username := c.Get("username").(string)
	sortBy := c.QueryParam("sortBy")

	reversed, err := strconv.ParseBool(c.QueryParam("reversed"))
	if err != nil {
		reversed = false
	}

	var urls []domain.Url
	if strings.TrimSpace(sortBy) != "" {
		urls, err = h.UrlGetter.GetAllSorted(context.Background(), username, sortBy, reversed)
	} else {
		urls, err = h.UrlGetter.GetAll(context.Background(), username)
	}

	if err != nil {
		return c.Render(http.StatusOK, "components/error", "Failed too load urls from the database")
	}

	data := components.UrlsTableData{
		Urls:     urls,
		SortedBy: sortBy,
		Reversed: !reversed,
	}
	return c.Render(http.StatusOK, "components/urls-table", data)
}

func (h UrlHandler) Listen(c echo.Context) error {
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")

	username := c.Param("username")
	hash := c.Param("hash")

	url, err := h.UrlGetter.Get(context.Background(), username, hash)
	if err != nil {
		return c.String(http.StatusOK, "-")
	}

	for url.ExpiresAt.After(time.Now().In(time.UTC)) {
		expiresIn := time.Until(url.ExpiresAt).Round(time.Second)
		sse.Send(c.Response().Writer, "url_update", []byte(expiresIn.String()))
		c.Response().Flush()
		time.Sleep(time.Second)
	}

	sse.Send(c.Response().Writer, "url_expired", []byte("Expired"))
	c.Response().Flush()
	return c.NoContent(http.StatusOK)
}

func (h UrlHandler) Create(c echo.Context) error {
	origin := c.FormValue("origin")
	origin = strings.TrimSuffix(origin, "/")

	if err := validation.ValidateUrl(origin); err != nil {
		return c.Render(http.StatusOK, "components/error", err.Error())
	}

	expiresIn, err := strconv.ParseInt(c.FormValue("expires_in"), 10, 8)
	if err != nil {
		return c.Render(http.StatusOK, "components/error", "Failed to parse expiration time")
	}

	if err := validation.ValidateExpiration(int(expiresIn)); err != nil {
		return c.Render(http.StatusOK, "components/error", err.Error())
	}

	url := domain.Url{
		Username:  c.Get("username").(string),
		Hash:      fmt.Sprintf("%08x", crc32.ChecksumIEEE([]byte(origin))),
		Origin:    origin,
		ExpiresAt: time.Now().In(time.UTC).Add(time.Duration(expiresIn) * time.Minute),
	}

	if err := h.UrlCreator.Create(context.Background(), url); err != nil {
		return c.Render(http.StatusOK, "components/error", "Failed to save the url in the database")
	}

	return c.NoContent(http.StatusOK)
}

func (h UrlHandler) Update(c echo.Context) error {
	username := c.Param("username")
	hash := c.Param("hash")
	newOrigin := c.FormValue("new-origin")

	if err := validation.ValidateUrl(newOrigin); err != nil {
		return c.Render(http.StatusOK, "components/error", err.Error())
	}

	url, err := h.UrlGetter.Get(context.Background(), username, hash)
	if err != nil {
		return c.Render(http.StatusOK, "components/error", "Failed to find link in the database")
	}

	updates := domain.UrlUpdates{
		Origin:     newOrigin,
		Hash:       fmt.Sprintf("%08x", crc32.ChecksumIEEE([]byte(newOrigin))),
		LastUsedAt: url.LastUsedAt,
	}
	if err := h.UrlUpdater.Update(context.Background(), username, hash, updates); err != nil {
		return c.Render(http.StatusOK, "components/error", "Failed to update the origin")
	}

	c.Response().Header().Set("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}

func (h UrlHandler) Delete(c echo.Context) error {
	username := c.Param("username")
	hash := c.Param("hash")

	if err := h.UrlDeleter.Delete(context.Background(), username, hash); err != nil {
		return c.NoContent(http.StatusNoContent)
	}

	return c.NoContent(http.StatusOK)
}

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
	UrlService domain.UrlService
}

func (h UrlHandler) GetOriginDialog(c echo.Context) error {
	body := make(map[string]any)
	if err := c.Bind(&body); err != nil {
		return c.Render(http.StatusOK, "components/error", "Invalid request")
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

func (h UrlHandler) GetOrigin(c echo.Context) error {
	c.Response().Header().Set("Cache-Control", "no-cache, max-age=0")

	username := c.Param("username")
	hash := c.Param("hash")

	url, err := h.UrlService.Get(context.Background(), username, hash)
	if err != nil {
		return c.Render(http.StatusOK, "404", nil)
	}

	if url.ExpiresAt.Before(time.Now().In(time.UTC)) {
		return c.Render(http.StatusOK, "404", nil)
	}

	updates := domain.UrlUpdates{
		Origin:     url.Origin,
		Hash:       url.Hash,
		LastUsedAt: time.Now(),
	}
	h.UrlService.Update(context.Background(), username, hash, updates)
	return c.Redirect(http.StatusMovedPermanently, url.Origin)
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
		urls, err = h.UrlService.GetAllSorted(context.Background(), username, sortBy, reversed)
	} else {
		urls, err = h.UrlService.GetAll(context.Background(), username)
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

	url, err := h.UrlService.Get(context.Background(), username, hash)
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
	return c.NoContent(http.StatusOK)
}

func (h UrlHandler) Create(c echo.Context) error {
	origin := c.FormValue("origin")
	origin = strings.TrimSuffix(origin, "/")

	url := domain.Url{
		Username:  c.Get("username").(string),
		Hash:      fmt.Sprintf("%08x", crc32.ChecksumIEEE([]byte(origin))),
		Origin:    origin,
		ExpiresAt: time.Now().In(time.UTC).Add(5 * time.Second),
	}

	if err := validation.ValidateUrl(origin); err != nil {
		return c.Render(http.StatusOK, "components/error", err.Error())
	}

	if err := h.UrlService.Create(context.Background(), url); err != nil {
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

	url, err := h.UrlService.Get(context.Background(), username, hash)
	if err != nil {
		return c.Render(http.StatusOK, "components/error", "Failed to find link in the database")
	}

	updates := domain.UrlUpdates{
		Origin:     newOrigin,
		Hash:       fmt.Sprintf("%08x", crc32.ChecksumIEEE([]byte(newOrigin))),
		LastUsedAt: url.LastUsedAt,
	}
	if err := h.UrlService.Update(context.Background(), username, hash, updates); err != nil {
		return c.Render(http.StatusOK, "components/error", "Failed to update the origin")
	}

	c.Response().Header().Set("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}

func (h UrlHandler) Delete(c echo.Context) error {
	username := c.Param("username")
	hash := c.Param("hash")

	if err := h.UrlService.Delete(context.Background(), username, hash); err != nil {
		return c.NoContent(http.StatusNoContent)
	}

	return c.NoContent(http.StatusOK)
}

package handlers

import (
	"context"
	"fmt"
	"hash/crc32"
	"net/http"
	"strings"
	"time"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/SergeyCherepiuk/surl/pkg/http/validation"
	"github.com/SergeyCherepiuk/surl/public/views/components"
	"github.com/labstack/echo/v4"
)

type UrlHandler struct {
	UrlService domain.UrlService
}

func (h UrlHandler) GetOrigin(c echo.Context) error {
	c.Response().Header().Set("Cache-Control", "no-cache, max-age=0")

	username := c.Param("username")
	hash := c.Param("hash")

	url, err := h.UrlService.Get(context.Background(), username, hash)
	if err != nil {
		return c.Render(http.StatusOK, "404", nil)
	}

	updates := domain.UrlUpdates{
		Origin:     url.Origin,
		LastUsedAt: time.Now(),
	}
	h.UrlService.Update(context.Background(), username, hash, updates)
	return c.Redirect(http.StatusMovedPermanently, url.Origin)
}

func (h UrlHandler) GetAll(c echo.Context) error {
	username := c.Get("username").(string)
	sortBy := c.QueryParam("sortBy")

	var urls []domain.Url
	var err error
	if strings.TrimSpace(sortBy) != "" {
		urls, err = h.UrlService.GetAllSorted(context.Background(), username, sortBy)
	} else {
		urls, err = h.UrlService.GetAll(context.Background(), username)
	}

	if err != nil {
		return c.String(http.StatusOK, "Failed too load urls from the database")
	}

	data := components.UrlsTableComponentData{Urls: urls}
	return c.Render(http.StatusOK, "components/urls-table", data)
}

func (h UrlHandler) Create(c echo.Context) error {
	origin := c.FormValue("origin")
	origin = strings.TrimSuffix(origin, "/")

	url := domain.Url{
		Username: c.Get("username").(string),
		Hash:     fmt.Sprintf("%08x", crc32.ChecksumIEEE([]byte(origin))),
		Origin:   origin,
	}

	if err := validation.ValidateOrigin(origin); err != nil {
		return c.String(http.StatusOK, err.Error())
	}

	if err := h.UrlService.Create(context.Background(), url); err != nil {
		return c.String(http.StatusOK, "Failed to save the url in the database")
	}

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

package handlers

import (
	"context"
	"fmt"
	"hash/crc32"
	"net/http"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/SergeyCherepiuk/surl/pkg/http/validation"
	"github.com/labstack/echo/v4"
)

type UrlHandler struct {
	UrlService domain.UrlService
}

func (h UrlHandler) GetAll(c echo.Context) error {
	username := c.Get("username").(string)

	urls, err := h.UrlService.GetAll(context.Background(), username)
	if err != nil {
		return c.String(http.StatusOK, "Failed too load urls from the database")
	}

	return c.Render(http.StatusOK, "components/urls-table-content", urls)
}

func (h UrlHandler) Create(c echo.Context) error {
	username := c.Get("username").(string)
	origin := c.FormValue("origin")
	hash := crc32.ChecksumIEEE([]byte(origin))

	url := domain.Url{
		Username: username,
		Hash:     fmt.Sprintf("%08x", hash),
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

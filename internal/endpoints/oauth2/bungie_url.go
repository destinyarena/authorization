package oauth2

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) bungieURL(c echo.Context) error {
	oauthurl, err := h.Providers.Bungie.FetchURL(nil)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, oauthurl)
}

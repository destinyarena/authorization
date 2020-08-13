package oauth2

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) faceitURL(c echo.Context) error {
	oauthurl, err := h.Providers.Faceit.FetchURL(nil)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// FACEIT IS A DOGSHIT PLATFORM THAT DOESN'T RESPECT THE OAUTH SPEC REEEE
	oauthurl = fmt.Sprintf("%s&redirect_popup=true", oauthurl)

	return c.String(http.StatusOK, oauthurl)
}

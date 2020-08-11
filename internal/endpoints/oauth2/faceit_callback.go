package oauth2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/destinyarena/authorization/pkg/oauth2"
	"github.com/labstack/echo/v4"
)

type faceitUser struct {
	GUID     string `json:"guid"`
	Nickname string `json:"nickname"`
}

func (h *handler) getFaceitUser(token string) (*faceitUser, error) {
	client := new(http.Client)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/auth/v1/resources/userinfo", h.Providers.Faceit.GetBaseAPI()), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("Error Code: %d", res.StatusCode)
	}

	var user faceitUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (h *handler) faceitCallback(c echo.Context) error {
	payload := new(OauthPayload)
	if err := c.Bind(payload); err != nil {
		return c.String(http.StatusBadRequest, "Invalid or missing code")
	}

	tokenpayload := oauth2.Token{}
	if err := h.Providers.Faceit.GetToken(payload.Code, &tokenpayload); err != nil {
		h.Logger.Error(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	user, err := h.getFaceitUser(tokenpayload.AccessToken)
	if err != nil {
		h.Logger.Error(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// TODO add JWT encapsulation and return that instead

	return c.JSON(http.StatusOK, user)
}

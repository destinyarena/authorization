package oauth2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/destinyarena/authorization/pkg/oauth2"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type (
	faceitUser struct {
		GUID     string `json:"guid"`
		Nickname string `json:"nickname"`
	}

	jwtFaceitClaim struct {
		faceitUser
		jwt.StandardClaims
	}
)

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

	claims := &jwtFaceitClaim{
		faceitUser: *user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(1)).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token, err := h.JWTManager.Sign(claims)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	data := map[string]interface{}{
		"token": token,
		"user":  &user,
	}

	return c.JSON(http.StatusOK, data)
}

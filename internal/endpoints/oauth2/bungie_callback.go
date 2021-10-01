package oauth2

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/destinyarena/authorization/pkg/oauth2"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type (
	bungieUserResponse struct {
		Response    *bungieUserMembership `json:"Response" validate:"required"`
		Message     string                `json:"Message"`
		ErrorStatus string                `json:"ErrorStatus"`
		ErrorCode   int                   `json:"ErrorCode"`
	}

	bungieUserMembership struct {
		BungieNetUser *bungieUser `json:"bungieNetUser" validate:"required"`
	}

	bungieUser struct {
		ID                  string     `json:"membershipId" validate:"required"`
		DisplayName         string     `json:"displayName" validate:"required"`
		SteamDisplayName    string     `json:"steamDisplayName,omitempty"`
		XboxDisplayName     string     `json:"xboxDisplayName,omitempty"`
		PSNDisplayName      string     `json:"psnDisplayName,omitempty"`
		BlizzardDisplayName string     `json:"blizzardDisplayName,omitempty"`
		TwitchDisplayName   string     `json:"twitchDisplayName,omitempty"`
		FirstAccess         *time.Time `json:"firstAccess,omitempty"`
	}

	jwtBungieClaim struct {
		bungieUser
		jwt.StandardClaims
	}
)

func (h *handler) getBungieUser(token string) (*bungieUser, error) {
	client := new(http.Client)

	requrl := fmt.Sprintf("%s/User/GetMembershipsForCurrentUser", h.Providers.Bungie.GetBaseAPI())
	req, err := http.NewRequest("GET", requrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("X-API-Key", h.Config.BungieAPIKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("Error code: %d", resp.StatusCode)
	}

	var payload bungieUserResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	return payload.Response.BungieNetUser, nil
}

func (h *handler) bungieCallback(c echo.Context) error {
	payload := new(OauthPayload)
	if err := c.Bind(payload); err != nil {
		h.Logger.Error(err)
		return c.String(http.StatusBadRequest, "Error while processing payload")
	}

	tokenpayload := oauth2.Token{}
	if err := h.Providers.Bungie.GetToken(payload.Code, &tokenpayload); err != nil {
		h.Logger.Error(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	user, err := h.getBungieUser(tokenpayload.AccessToken)
	if err != nil {
		h.Logger.Error(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	/* SteamID Check
	if len(user.SteamDisplayName) == 0 {
		err = errors.New("You must have a steam account linked")
		h.Logger.Error(err)
		return c.String(401, err.Error())
	}
	*/

	if user.FirstAccess == nil {
		err = errors.New("Looks like you've never played Destiny 2 before")
		h.Logger.Error(err)
		return c.String(401, err.Error())
	}

	bestbefore := time.Now().Add(-730 * time.Hour)

	if !bestbefore.After(*user.FirstAccess) {
		h.Logger.Infof("Account ID: %s Name: %s Created at: %v", user.ID, user.DisplayName, (*user.FirstAccess))
		return c.String(401, "Your account must be older than 30 days to play faceit.")
	}

	claims := &jwtBungieClaim{
		bungieUser: *user,
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

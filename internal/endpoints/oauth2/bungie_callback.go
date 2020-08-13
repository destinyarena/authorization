package oauth2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/destinyarena/authorization/pkg/oauth2"
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
		//PrimaryMembershipID int         `json:"primaryMembershipId"`
		BungieNetUser *bungieUser `json:"bungieNetUser" validate:"required"`
	}

	bungieUser struct {
		ID                  string `json:"membershipId" validate:"required"`
		DisplayName         string `json:"displayName" validate:"required"`
		SteamDisplayName    string `json:"steamDisplayName,omitempty"`
		XboxDisplayName     string `json:"xboxDisplayName,omitempty"`
		PSNDisplayName      string `json:"psnDisplayName,omitempty"`
		BlizzardDisplayName string `json:"blizzardDisplayName,omitempty"`
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

	// TODO add JWT encapsulation and return that instead of the user struct

	return c.JSON(http.StatusOK, user)

}

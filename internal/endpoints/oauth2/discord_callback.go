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
	discordUser struct {
		ID            string `json:"id"`
		Username      string `json:"username"`
		Discriminator string `json:"discriminator"`
	}

	discordGuild struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
)

func (h *handler) checkDiscordGuilds(token string) (bool, error) {
	client := new(http.Client)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/@me/guilds", h.Providers.Discord.GetBaseAPI()), nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := client.Do(req)
	if err != nil {
		return false, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, nil
	}

	if res.StatusCode > 299 {
		return false, fmt.Errorf("Status Code: %d", res.StatusCode)
	}

	guilds := make([]discordGuild, 0)
	if err := json.Unmarshal(body, &guilds); err != nil {
		return false, err
	}

	for _, guild := range guilds {
		if guild.ID == h.Config.DiscordGuildID {
			return true, nil
		}
	}

	return false, nil

}

func (h *handler) getDiscordUser(token string) (*discordUser, error) {
	client := new(http.Client)

	requrl := fmt.Sprintf("%s/users/@me", h.Providers.Discord.GetBaseAPI())
	req, err := http.NewRequest("GET", requrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

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
		return nil, fmt.Errorf("Error code: %d", res.StatusCode)
	}

	var user discordUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (h *handler) discordCallback(c echo.Context) error {
	payload := new(OauthPayload)
	if err := c.Bind(payload); err != nil {
		h.Logger.Error(err)
		return c.String(http.StatusBadRequest, "Error while processing payload")
	}

	tokenpayload := oauth2.Token{}
	if err := h.Providers.Discord.GetToken(payload.Code, &tokenpayload); err != nil {
		h.Logger.Error(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	user, err := h.getDiscordUser(tokenpayload.AccessToken)
	if err != nil {
		h.Logger.Error(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	ok, err := h.checkDiscordGuilds(tokenpayload.AccessToken)
	if !ok {
		return c.String(http.StatusUnauthorized, "please join our server before attempting to register")
	} else if err != nil {
		h.Logger.Error(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// TODO add JWT encapsulation and return that instead of user json

	return c.JSON(http.StatusOK, user)
}

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

func (h *handler) getDiscordGuilds(token string) ([]*discordGuild, error) {
	client := new(http.Client)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/@me/guilds", h.Providers.Discord.GetBaseAPI()), nil)
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
		return nil, fmt.Errorf("Status Code: %d", res.StatusCode)
	}

	guilds := make([]*discordGuild, 0)
	if err := json.Unmarshal(body, &guilds); err != nil {
		return nil, err
	}

	return guilds, nil
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

	guilds, err := h.getDiscordGuilds(tokenpayload.AccessToken)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if ok := checkDiscordGuilds(guilds, h.Config.DiscordGuildID); !ok {
		return c.String(http.StatusUnauthorized, "please join our server before attempting to register")
	}

	if banned := bannedDiscordGuilds(guilds, h.Config.DiscordBannedGuilds); banned {
		return c.String(http.StatusForbidden, "Looks like you are in a banned guild")
	}

	// TODO add JWT encapsulation and return that instead of user json

	return c.JSON(http.StatusOK, user)
}

func checkDiscordGuilds(guilds []*discordGuild, ID string) bool {
	for _, guild := range guilds {
		if guild.ID == ID {
			return true
		}
	}

	return false
}

func bannedDiscordGuilds(guilds []*discordGuild, banned []string) bool {
	for _, guild := range guilds {
		for _, bguild := range banned {
			if guild.ID == bguild {
				return true
			}
		}
	}

	return false
}

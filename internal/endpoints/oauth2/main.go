package oauth2

import (
	"github.com/destinyarena/authorization/internal/jwtmanager"
	"github.com/destinyarena/authorization/pkg/oauth2"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type (
	providers struct {
		Faceit  oauth2.Oauth2
		Discord oauth2.Oauth2
		Bungie  oauth2.Oauth2
	}
	handler struct {
		Logger     *logrus.Logger
		Config     *Config
		Providers  *providers
		JWTManager jwtmanager.JWTManager
	}

	// Handler exports the register hook used to register with echo
	Handler interface {
		Register(*echo.Group)
	}
)

// New creates a new handler
func New(logger *logrus.Logger, config *Config, jwtManager jwtmanager.JWTManager) (Handler, error) {
	discordProvider := &oauth2.Provider{
		AuthorizeURL: "https://discord.com/oauth2/authorize",
		BaseAPI:      "https://discord.com/api/v6",
		TokenURI:     "/oauth2/token",
		ClientID:     config.DiscordClientID,
		ClientSecret: config.DiscordClientSecret,
		ResponseType: "code",
		Scope:        "identify guilds",
		RedirectURI:  &config.DiscordRedirectURI,
	}

	discord, err := oauth2.New(discordProvider)
	if err != nil {
		return nil, err
	}

	faceitProvider := &oauth2.Provider{
		AuthorizeURL: "https://cdn.faceit.com/widgets/sso/index.html",
		BaseAPI:      "https://api.faceit.com",
		TokenURI:     "/auth/v1/oauth/token",
		ClientID:     config.FaceitClientID,
		ClientSecret: config.FaceitClientSecret,
		ResponseType: "code",
		Scope:        "",
	}

	faceit, err := oauth2.New(faceitProvider)
	if err != nil {
		return nil, err
	}

	bungieProvider := &oauth2.Provider{
		AuthorizeURL: "https://www.bungie.net/en/OAuth/Authorize",
		BaseAPI:      "https://www.bungie.net/Platform",
		TokenURI:     "/app/oauth/token",
		ClientID:     config.BungieClientID,
		ClientSecret: config.BungieClientSecret,
		ResponseType: "code",
		Scope:        "",
		//Scope:        "ReadBasicUserProfile",
	}

	bungie, err := oauth2.New(bungieProvider)
	if err != nil {
		return nil, err
	}

	providers := &providers{
		Discord: discord,
		Bungie:  bungie,
		Faceit:  faceit,
	}

	return &handler{
		Logger:     logger,
		Config:     config,
		Providers:  providers,
		JWTManager: jwtManager,
	}, nil
}

func (h *handler) Register(g *echo.Group) {
	g.GET("/discord/url", h.discordURL)
	g.GET("/discord/callback", h.discordCallback)

	g.GET("/bungie/url", h.bungieURL)
	g.GET("/bungie/callback", h.bungieCallback)

	g.GET("/faceit/url", h.faceitURL)
	g.GET("/faceit/callback", h.faceitCallback)

}

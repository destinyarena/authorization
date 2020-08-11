package oauth2

import "github.com/destinyarena/authorization/internal/utils"

type (
	// Config is the oauth2 configuration
	Config struct {
		DiscordGuildID      string   `env:"DISCORD_GUILD_ID,required"`
		DiscordClientID     string   `env:"DISCORD_CLIENT_ID,required"`
		DiscordClientSecret string   `env:"DISCORD_CLIENT_SECRET,required"`
		DiscordRedirectURI  string   `env:"DISCORD_REDIRECT_URI,required"`
		DiscordBannedGuilds []string `env:"DISCORD_BANNED_GUILDS" envSperator:","`

		BungieClientID     string `env:"BUNGIE_CLIENT_ID,required"`
		BungieClientSecret string `env:"BUNGIE_CLIENT_SECRET,required"`
		BungieAPIKey       string `env:"BUNGIE_API_KEY,required"`

		FaceitClientID     string `env:"FACEIT_CLIENT_ID,required"`
		FaceitClientSecret string `env:"FACEIT_CLIENT_SECRET,required"`
	}
)

// LoadConfig returns config for oauth2
func LoadConfig() (*Config, error) {
	cfg := Config{}
	if err := utils.EnvLoader(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

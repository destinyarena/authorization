package oauth2

type (
	// Config is the oauth2 configuration
	Config struct {
		DiscordGuildID      string `env:"DISCORD_GUILD_ID,required"`
		DiscordClientID     string `env:"DISCORD_CLIENT_ID,required"`
		DiscordClientSecret string `env:"DISCORD_CLIENT_SECRET,required"`
		DiscordRedirectURI  string `env:"DISCORD_REDIRECT_URI,required"`

		BungieClientID     string `env:"BUNGIE_CLIENT_ID,required"`
		BungieClientSecret string `env:"BUNGIE_CLIENT_SECRET,required"`
		BungieAPIKey       string `env:"BUNGIE_API_KEY,required"`
		BungieRedirectURI  string `env:"BUNGIE_REDIRECT_URI,required"`

		FaceitClientID     string `env:"FACEIT_CLIENT_ID,required"`
		FaceitClientSecret string `env:"FACEIT_CLIENT_SECRET,required"`
		FaceitRedirectURI  string `env:"FACEIT_REDIRECT_URI,required"`
	}
)

// LoadConfig returns config for oauth2
func LoadConfig() (*Config, error) {

	return &Config{}, nil
}

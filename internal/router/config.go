package router

type (
	// Config is the configuration for echo router
	Config struct {
		Host string `env:"HOST"`
		Port string `env:"PORT"`
	}
)

// LoadConfig loads the echo router config
func LoadConfig() (*Config, error) {
	return &Config{
		Host: "0.0.0.0",
		Port: "3333",
	}, nil
}

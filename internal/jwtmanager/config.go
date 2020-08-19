package jwtmanager

import "github.com/destinyarena/authorization/internal/utils"

type (
	// Config holds configuration info for jwtmanager
	Config struct {
		Secret string `env:"JWT_SECRET" envDefault:"testing"`
	}
)

func newEnvConfig() (*Config, error) {
	cfg := new(Config)
	if err := utils.EnvLoader(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

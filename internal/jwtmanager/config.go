package jwtmanager

import "github.com/destinyarena/authorization/internal/utils"

type (
	// Config holds configuration info for jwtmanager
	Config struct {
		PrivKeyPath string `env:"JWT_PRIV_KEY_PATH" envDefault:"/etc/keys/rsa"`
		PubKeyPath  string `env:"JWT_PUB_KEY_PATH" envDefault:"/etc/keys/rsa.pub"`
	}
)

func newEnvConfig() (*Config, error) {
	cfg := new(Config)
	if err := utils.EnvLoader(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

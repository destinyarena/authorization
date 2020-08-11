package utils

import "github.com/caarlos0/env"

// EnvLoader populates a struct with env values
func EnvLoader(cfg interface{}) error {
	return env.Parse(cfg)
}

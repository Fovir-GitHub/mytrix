package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

func Load() error {
	Config = &MytrixConfig{}

	if err := env.ParseWithOptions(Config, env.Options{
		Prefix: "MYTRIX_",
	}); err != nil {
		return fmt.Errorf("load config error: %w", err)
	}
	return nil
}

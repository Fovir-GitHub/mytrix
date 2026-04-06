package config

import (
	"fmt"
	"log/slog"
	"time"
	_ "time/tzdata"

	"github.com/caarlos0/env/v11"
)

func Load() error {
	Config = &MytrixConfig{}

	if err := env.ParseWithOptions(Config, env.Options{
		Prefix: "MYTRIX_",
	}); err != nil {
		return fmt.Errorf("load config error: %w", err)
	}

	return Config.validate()
}

func SetTimeZone() {
	tz := Config.TZ
	if tz != "" {
		loc, err := time.LoadLocation(tz)
		if err != nil {
			slog.Warn(
				"invalid timezone, use default location",
				"tz", tz,
				"default", time.Local.String(),
				"err", err,
			)
		} else {
			time.Local = loc
		}
	}
	slog.Info("set timezone", "timezone", time.Local.String())
}

package config

import (
	"fmt"
	"log/slog"
	"time"
	_ "time/tzdata"

	"github.com/caarlos0/env/v11"
)

// Load loads the configuration from environment variables, returning an error if unsuccessful.
func Load() error {
	Config = &MytrixConfig{}

	if err := env.ParseWithOptions(Config, env.Options{
		Prefix: "MYTRIX_",
	}); err != nil {
		return fmt.Errorf("load config error: %w", err)
	}

	return Config.validate()
}

// SetTimeZone sets the bot's timezone from the environment variable.
// It loads the specified timezone and assigns it to time.Local, falling back to the
// local timezone if the specified timezone is invalid.
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

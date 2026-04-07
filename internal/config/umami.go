package config

import (
	"fmt"
	"slices"
)

type UmamiConfig struct {
	Enabled         bool   `env:"UMAMI_ENABLED" envDefault:"false"`
	Server          string `env:"UMAMI_SERVER"`
	Username        string `env:"UMAMI_USERNAME"`
	Password        string `env:"UMAMI_PASSWORD"`
	DefaultInterval string `env:"UMAMI_DEFAULT_INTERVAL" envDefault:"daily"`
	Format          string `env:"UMAMI_FORMAT" envDefault:"- {{.Name}} - {{.Domain}}\n\t- Visitors: {{.Visitors}}\n\t- Visits: {{.Visits}}\n\tBounces Rate: {{.BouncesRate}}"`
}

func (mc *MytrixConfig) validateUmami() error {
	cfg := mc.Umami
	if !cfg.Enabled {
		return nil
	}
	if cfg.Server == "" || cfg.Username == "" || cfg.Password == "" {
		return fmt.Errorf("MYTRIX_UMAMI_SERVER, MYTRIX_UMAMI_USERNAME and MYTRIX_UMAMI_UMAMI_PASSWORD are required when MYTRIX_UMAMI_ENABLED=true")
	}
	if !validUmamiInterval(cfg.DefaultInterval) {
		return fmt.Errorf("invalid umami default interval: %s", cfg.DefaultInterval)
	}
	return nil
}

func validUmamiInterval(interval string) bool {
	validIntervals := []string{"daily", "weekly", "monthly", "yearly"}
	return slices.Contains(validIntervals, interval)
}

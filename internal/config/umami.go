package config

import "fmt"

type UmamiConfig struct {
	Enabled  bool   `env:"UMAMI_ENABLED" envDefault:"false"`
	Server   string `env:"UMAMI_SERVER"`
	Username string `env:"UMAMI_USERNAME"`
	Password string `env:"UMAMI_PASSWORD"`
	Format   string `env:"UMAMI_FORMAT" envDefault:"- {{.Name}} - {{.Domain}}\n\t- Visitors: {{.Visitors}}\n\t- Visits: {{.Visits}}\n\tBounces Rate: {{.BouncesRate}}"`
}

func (mc *MytrixConfig) validateUmami() error {
	cfg := mc.Umami
	if !cfg.Enabled {
		return nil
	}
	if cfg.Server == "" || cfg.Username == "" || cfg.Password == "" {
		return fmt.Errorf("MYTRIX_UMAMI_SERVER, MYTRIX_UMAMI_USERNAME and MYTRIX_UMAMI_UMAMI_PASSWORD are required when MYTRIX_UMAMI_ENABLED=true")
	}
	return nil
}

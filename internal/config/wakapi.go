package config

import "fmt"

type WakapiConfig struct {
	// Enabled determines whether to enable Wakapi feature.
	Enabled bool `env:"WAKAPI_ENABLED" envDefault:"false"`
	// Server sets the Wakapi server.
	Server string `env:"WAKAPI_SERVER"`
	// APIKey sets the api key to access Wakapi API.
	APIKey string `env:"WAKAPI_API_KEY"`
	// UserID sets the user id.
	UserID string `env:"WAKAPI_USER_ID" envDefault:"current"`
	// DefaultInterval sets the default interval to fetch wakapi analysis.
	DefaultInterval string `env:"WAKAPI_DEFAULT_INTERVAL" envDefault:"today"`
}

func (mc *MytrixConfig) validateWakapi() error {
	cfg := mc.Wakapi
	if cfg.Enabled && (cfg.Server == "" || cfg.APIKey == "") {
		return fmt.Errorf("MYTRIX_WAKAPI_SERVER and MYTRIX_WAKAPI_API_KEY are required when MYTRIX_WAKAPI_ENABLED=true")
	}
	return nil
}

package config

import "fmt"

// GotifyConfig holds the configuration for Gotify integration.
type GotifyConfig struct {
	// Enabled determines whether to enable Gotify feature.
	Enabled bool `env:"GOTIFY_ENABLED" envDefault:"false"`
	// Server sets the server of Gotify.
	Server string `env:"GOTIFY_SERVER"`
	// Token is used to access Gotify.
	Token string `env:"GOTIFY_TOKEN"`
	// Format is the message style of Gotify (support Markdown).
	Format string `env:"GOTIFY_FORMAT" envDefault:"# {{.Title}}\n\n**{{.Message}}**\n\n- ID: {{.ID}}\n- Date: {{.Date}}"`
}

// validateGotify validates the Gotify configuration and ensures required fields are present when enabled.
func (mc *MytrixConfig) validateGotify() error {
	cfg := mc.Gotify
	if cfg.Enabled && (cfg.Server == "" || cfg.Token == "") {
		return fmt.Errorf("MYTRIX_GOTIFY_SERVER and MYTRIX_GOTIFY_TOKEN are required when MYTRIX_GOTIFY_ENABLE=true")
	}

	return nil
}

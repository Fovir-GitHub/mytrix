package config

type GotifyConfig struct {
	// Enabled determines whether to enable gotify feature.
	Enabled bool `env:"GOTIFY_ENABLED" envDefault:"false"`
	// Server sets the server of gotify.
	Server string `env:"GOTIFY_SERVER"`
	// Token is used to access gotify.
	Token string `env:"GOTIFY_TOKEN"`
}

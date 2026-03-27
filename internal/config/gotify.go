package config

type GotifyConfig struct {
	// Enabled determines whether to enable gotify feature.
	Enabled bool `env:"GOTIFY_ENABLED" envDefault:"false"`
	// Server sets the server of gotify.
	Server string `env:"GOTIFY_SERVER"`
	// Token is used to access gotify.
	Token string `env:"GOTIFY_TOKEN"`
	// Format is the message style of gotify (support Markdown).
	Format string `env:"GOTIFY_FORMAT" envDefault:"# {{.Title}}\n\n**{{.Message}}**\n\n- ID: {{.ID}}\n- Date: {{.Date}}"`
}

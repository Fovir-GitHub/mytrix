package config

type RSSConfig struct {
	// Enabled determines whether to enable RSS integration.
	Enabled bool `env:"RSS_ENABLED" envDefault:"false"`
	// Format sets the output format of RSS items.
	Format string `env:"RSS_FORMAT" envDefault:"{{.Title}} - {{.Link}}"`
}

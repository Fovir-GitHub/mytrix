package config

type RSSConfig struct {
	Enabled bool `env:"RSS_ENABLED" envDefault:"false"`
}

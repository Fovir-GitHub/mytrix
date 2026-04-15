package config

type RSSConfig struct {
	// Enabled determines whether to enable RSS integration.
	Enabled bool `env:"RSS_ENABLED" envDefault:"false"`
	// Format sets the output format of RSS items.
	Format string `env:"RSS_FORMAT" envDefault:"{{.Title}} - {{.Link}}"`
	// Cron sets the fetch interval of RSS feeds (hourly by default).
	Cron string `env:"RSS_CRON" envDefault:"0 * * * *"`
}

func (mc *MytrixConfig) validateRSS() error {
	cfg := mc.RSS
	if !cfg.Enabled {
		return nil
	}
	return mc.validateCrons([]string{cfg.Cron})
}

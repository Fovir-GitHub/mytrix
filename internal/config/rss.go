package config

// RSSConfig holds the configuration for RSS feed integration.
type RSSConfig struct {
	// Enabled determines whether to enable RSS integration.
	Enabled bool `env:"RSS_ENABLED" envDefault:"false"`
	// FeedFormat sets the output format of RSS feeds.
	FeedFormat string `env:"RSS_FEED_FORMAT" envDefault:"- {{.ID}} {{.Title}}: {{.URL}}"`
	// ItemFormat sets the output format of RSS items.
	ItemFormat string `env:"RSS_ITEM_FORMAT" envDefault:"{{.Title}} - {{.Link}}"`
	// Cron sets the fetch interval of RSS feeds (hourly by default).
	Cron string `env:"RSS_CRON" envDefault:"0 * * * *"`
	// UpdateAfterAdd determines whether to run `Update` after adding a feed.
	UpdateAfterAdd bool `env:"RSS_UPDATE_AFTER_ADD" envDefault:"true"`
}

// validateRSS validates the RSS configuration and ensures cron expressions are valid when enabled.
func (mc *MytrixConfig) validateRSS() error {
	cfg := mc.RSS
	if !cfg.Enabled {
		return nil
	}
	return mc.validateCrons([]string{cfg.Cron})
}

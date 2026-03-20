package config

var Config *MytrixConfig

type MytrixConfig struct {
	LogLevel   string `env:"LOG_LEVEL" envDefault:"INFO"`
	Homeserver string `env:"HOMESERVER,required"`
	Bot        BotConfig
}

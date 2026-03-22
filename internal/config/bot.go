package config

type BotConfig struct {
	Username    string `env:"BOT_USERNAME,required"`
	Password    string `env:"BOT_PASSWORD,required"`
	RecoveryKey string `env:"BOT_RECOVERY_KEY,required"`
	PickleKey   string `env:"BOT_PICKLE_KEY,required"`
}

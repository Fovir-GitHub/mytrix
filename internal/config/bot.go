package config

type BotConfig struct {
	UserID      string `env:"BOT_USER_ID,required"`
	AccessToken string `env:"BOT_ACCESS_TOKEN,required"`
	SessionID   string `env:"SESSION_ID"`
}

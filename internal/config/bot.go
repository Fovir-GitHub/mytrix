package config

// BotConfig contains configuration specific to the Matrix bot.
// All fields are required and loaded from environment variables with the BOT_ prefix.
type BotConfig struct {
	// Username is the Matrix user ID for the bot (required).
	Username string `env:"BOT_USERNAME,required"`
	// Password is the password for the bot account (required).
	Password string `env:"BOT_PASSWORD,required"`
	// RecoveryKey is the encryption recovery key for the bot (required).
	RecoveryKey string `env:"BOT_RECOVERY_KEY,required"`
	// PickleKey is used to encrypt the crypto storage (required).
	PickleKey string `env:"BOT_PICKLE_KEY,required"`
}

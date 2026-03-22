package config

// Config holds the global configuration instance for the application.
var Config *MytrixConfig

// MytrixConfig represents the application configuration structure.
// It is populated from environment variables with the MYTRIX_ prefix.
type MytrixConfig struct {
	// LogLevel specifies the logging level (DEBUG, INFO, WARN, ERROR).
	LogLevel string `env:"LOG_LEVEL" envDefault:"INFO"`
	// Homeserver is the Matrix homeserver URL (required).
	Homeserver string `env:"HOMESERVER,required"`
	// Bot contains bot-specific configuration.
	Bot BotConfig
}

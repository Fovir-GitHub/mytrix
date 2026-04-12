/**
Package config provides mytrix configuration settings and validations.

All module-related configurations (e.g., Gotify, Umami, etc.), should provides a `validate` function to check whether the configuration is valid.

If a module has schedulers, it should provides a cron list, and use `validateCrons` to validate its cron expressions.
*/

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
	// RoomID is the chat room with bot.
	RoomID string `env:"ROOM_ID,required"`
	// Datadir sets the data directory.
	Datadir string `env:"DATA_DIR" envDefault:"/data"`
	// Timeout defines the timeout of http request.
	Timeout int `env:"TIMEOUT" envDefault:"10"`
	// TZ sets the timezone.
	TZ string `env:"TZ" envDefault:""`
	// Bot contains bot-specific configuration.
	Bot BotConfig
	// Gotify contains Gotify-specific configuration.
	Gotify GotifyConfig
	// WS specifics the websocket configuration.
	WS WSConfig
	// Msg contains message-specific configuration.
	Msg MsgConfig
	// Wakapi contains Wakapi-specific configurartion.
	Wakapi WakapiConfig
	// Umami contains Umami-specific configuration.
	Umami UmamiConfig
}

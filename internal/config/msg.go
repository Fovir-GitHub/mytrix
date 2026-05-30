package config

type MsgConfig struct {
	// AllowMarkdown determines whether to allow markdown content in messages.
	AllowMarkdown bool `env:"MSG_ALLOW_MARKDOWN" envDefault:"true"`
	// AllowHTML determines whether to allow HTML content in messages.
	AllowHTML bool `env:"MSG_ALLOW_HTML" envDefault:"false"`
	// DefaultMaxLength sets the default or fallback max length of a message.
	DefaultMaxLength int `env:"MSG_DEFAULT_MAX_LENGTH" envDefault:"4096"`
}

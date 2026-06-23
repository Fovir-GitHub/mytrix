package config

type DBConfig struct {
	// Path sets the database path.
	Path string `env:"DB_PATH" envDefault:"sqlite.db"`
	// BusyTimeout defines the timeout when db is busy (in ms).
	BusyTimeout int `env:"DB_BUSY_TIMEOUT" envDefault:"5000"`
}

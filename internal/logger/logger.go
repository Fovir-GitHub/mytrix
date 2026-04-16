package logger

import (
	"log/slog"
	"os"

	"codeberg.org/Fovir/mytrix/internal/config"
	"github.com/lmittmann/tint"
)

// Init initializes the global logger with the configured log level and output formatting.
// It uses the tint package for colorized terminal output.
func Init() {
	var level slog.Level
	if err := level.UnmarshalText([]byte(config.Config.LogLevel)); err != nil {
		level = slog.LevelInfo
	}

	handler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      level,
		TimeFormat: "15:04:05",
	})

	slog.SetDefault(slog.New(handler))
}

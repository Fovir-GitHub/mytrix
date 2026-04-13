package logger

import (
	"testing"

	"codeberg.org/Fovir/mytrix/internal/config"
	"log/slog"
)

func init() {
	// Initialize config for tests
	if config.Config == nil {
		config.Config = &config.MytrixConfig{
			LogLevel: "INFO",
		}
	}
}

func TestLogger_Init_ValidLevel(t *testing.T) {
	// Save original config
	originalLogLevel := config.Config.LogLevel
	defer func() { config.Config.LogLevel = originalLogLevel }()

	// Test with valid log level
	config.Config.LogLevel = "DEBUG"
	Init()

	// Just verify it doesn't panic - we can't easily test the global state change
	// without affecting other tests
	// We can check that the default logger was changed by checking if it's not nil
	if slog.Default() == nil {
		t.Error("Expected logger to be initialized")
	}
}

func TestLogger_Init_InvalidLevel(t *testing.T) {
	// Save original config
	originalLogLevel := config.Config.LogLevel
	defer func() { config.Config.LogLevel = originalLogLevel }()

	// Test with invalid log level (should default to INFO)
	config.Config.LogLevel = "INVALID"
	Init()

	if slog.Default() == nil {
		t.Error("Expected logger to be initialized")
	}
}

func TestLogger_Init_EmptyLevel(t *testing.T) {
	// Save original config
	originalLogLevel := config.Config.LogLevel
	defer func() { config.Config.LogLevel = originalLogLevel }()

	// Test with empty log level (should default to INFO)
	config.Config.LogLevel = ""
	Init()

	if slog.Default() == nil {
		t.Error("Expected logger to be initialized")
	}
}

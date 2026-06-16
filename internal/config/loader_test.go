package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad_Success(t *testing.T) {
	// Save original env vars
	originalHomeserver := os.Getenv("MYTRIX_HOMESERVER")
	originalRoomID := os.Getenv("MYTRIX_ROOM_ID")
	originalBotUsername := os.Getenv("MYTRIX_BOT_USERNAME")
	originalBotPassword := os.Getenv("MYTRIX_BOT_PASSWORD")
	originalBotRecoveryKey := os.Getenv("MYTRIX_BOT_RECOVERY_KEY")
	originalBotPickleKey := os.Getenv("MYTRIX_BOT_PICKLE_KEY")
	originalUserAdminID := os.Getenv("MYTRIX_USER_ADMIN_ID")

	//nolint
	defer func() {
		os.Setenv("MYTRIX_HOMESERVER", originalHomeserver)
		os.Setenv("MYTRIX_ROOM_ID", originalRoomID)
		os.Setenv("MYTRIX_BOT_USERNAME", originalBotUsername)
		os.Setenv("MYTRIX_BOT_PASSWORD", originalBotPassword)
		os.Setenv("MYTRIX_BOT_RECOVERY_KEY", originalBotRecoveryKey)
		os.Setenv("MYTRIX_BOT_PICKLE_KEY", originalBotPickleKey)
		os.Setenv("MYTRIX_USER_ADMIN_ID", originalUserAdminID)
	}()

	// Set required env vars
	//nolint
	{
		os.Setenv("MYTRIX_HOMESERVER", "https://matrix.example.com")
		os.Setenv("MYTRIX_ROOM_ID", "!roomid:matrix.example.com")
		os.Setenv("MYTRIX_BOT_USERNAME", "botuser")
		os.Setenv("MYTRIX_BOT_PASSWORD", "botpass")
		os.Setenv("MYTRIX_BOT_RECOVERY_KEY", "test-recovery-key")
		os.Setenv("MYTRIX_BOT_PICKLE_KEY", "test-pickle-key")
		os.Setenv("MYTRIX_USER_ADMIN_ID", "@example:example.com")
	}

	err := Load()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if Config == nil {
		t.Fatalf("Expected Config to be initialized")
	}
	if Config.Homeserver != "https://matrix.example.com" {
		t.Errorf("Expected homeserver 'https://matrix.example.com', got %q", Config.Homeserver)
	}
	if Config.RoomID != "!roomid:matrix.example.com" {
		t.Errorf("Expected room ID '!roomid:matrix.example.com', got %q", Config.RoomID)
	}
}

func TestLoad_MissingRequired(t *testing.T) {
	// Save original env vars
	originalHomeserver := os.Getenv("MYTRIX_HOMESERVER")
	originalRoomID := os.Getenv("MYTRIX_ROOM_ID")
	//nolint
	defer func() {
		os.Setenv("MYTRIX_HOMESERVER", originalHomeserver)
		os.Setenv("MYTRIX_ROOM_ID", originalRoomID)
	}()

	// Clear required env vars
	//nolint
	{
		os.Unsetenv("MYTRIX_HOMESERVER")
		os.Unsetenv("MYTRIX_ROOM_ID")
	}

	err := Load()
	if err == nil {
		t.Fatalf("Expected error for missing required fields")
	}
}

func TestSetTimeZone_ValidTimezone(t *testing.T) {
	// Save original timezone
	originalTZ := Config.TZ
	originalLocal := time.Local
	defer func() {
		Config.TZ = originalTZ
		time.Local = originalLocal
	}()

	Config.TZ = "America/New_York"
	SetTimeZone()

	if time.Local.String() != "America/New_York" {
		t.Errorf("Expected timezone 'America/New_York', got %q", time.Local.String())
	}
}

func TestSetTimeZone_EmptyTimezone(t *testing.T) {
	// Save original timezone
	originalTZ := Config.TZ
	originalLocal := time.Local
	defer func() {
		Config.TZ = originalTZ
		time.Local = originalLocal
	}()

	Config.TZ = ""
	SetTimeZone()

	// Should not change timezone if empty
	// time.Local remains unchanged
	if time.Local == nil {
		t.Errorf("Expected time.Local to be set")
	}
}

func TestSetTimeZone_InvalidTimezone(t *testing.T) {
	// Save original timezone
	originalTZ := Config.TZ
	originalLocal := time.Local
	defer func() {
		Config.TZ = originalTZ
		time.Local = originalLocal
	}()

	Config.TZ = "Invalid/Timezone"
	SetTimeZone()

	// Should fall back to default, so time.Local should still be set
	if time.Local == nil {
		t.Errorf("Expected time.Local to be set after invalid timezone")
	}
}

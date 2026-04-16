package storage

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"codeberg.org/Fovir/mytrix/internal/config"
	"maunium.net/go/mautrix"
)

// Session represents a Matrix client session containing user ID, access token, and device ID.
type Session struct {
	UserID      string `json:"user_id"`
	AccessToken string `json:"access_token"`
	DeviceID    string `json:"device_id"`
}

// sessionFile returns the path to the file where session data is stored.
func sessionFile() string {
	return filepath.Join(config.Config.Datadir, "session.json")
}

// LoadSession reads and decodes the session from the session file.
// It returns a pointer to the Session struct and any error encountered.
func LoadSession() (*Session, error) {
	var s Session
	data, err := os.ReadFile(sessionFile())
	slog.Debug(
		"load session",
		"session_file", sessionFile,
	)
	if err != nil {
		return nil, fmt.Errorf("load session failed: %w", err)
	}
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("unmarshal session file failed: %w", err)
	}
	return &s, nil
}

// SaveSession encodes and saves the session data from a login response to the session file.
// It takes a RespLogin object and extracts the session information to store.
func SaveSession(resp *mautrix.RespLogin) error {
	slog.Debug(
		"save session start",
		"session_file", sessionFile,
	)
	s := Session{
		UserID:      resp.UserID.String(),
		AccessToken: resp.AccessToken,
		DeviceID:    resp.DeviceID.String(),
	}

	data, _ := json.MarshalIndent(s, "", "  ")
	return os.WriteFile(sessionFile(), data, 0o600)
}

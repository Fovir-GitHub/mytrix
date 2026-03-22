package storage

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"maunium.net/go/mautrix"
)

type Session struct {
	UserID      string `json:"user_id"`
	AccessToken string `json:"access_token"`
	DeviceID    string `json:"device_id"`
}

const sessionFile = "db/session.json"

func LoadSession() (*Session, error) {
	var s Session
	data, err := os.ReadFile(sessionFile)
	slog.Debug(
		"load session",
		"session_file", sessionFile,
	)
	if err != nil {
		return nil, fmt.Errorf("load session failed: %w", err)
	}
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("unmarshal session file faile: %w", err)
	}
	return &s, nil
}

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
	return os.WriteFile(sessionFile, data, 0600)
}

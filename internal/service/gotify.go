package service

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"codeberg.org/Fovir/mytrix/internal/config"
	"codeberg.org/Fovir/mytrix/internal/model"
)

// GotifyService handles incoming Gotify webhook events.
// It provides functionality to process and decode Gotify notifications.
type GotifyService struct{}

// newGotifyService creates a GotifyService based on the configuration.
// It returns a new GotifyService instance if Gotify is enabled, otherwise it returns nil.
func newGotifyService() *GotifyService {
	if config.Config.Gotify.Enabled {
		slog.Info("gotify enabled")
		return &GotifyService{}
	}
	slog.Info("gotify disabled")
	return nil
}

func (g *GotifyService) HandleEvent(event *model.WsEvent) (*model.GotifyMessage, error) {
	var msg model.GotifyMessage
	data := event.Data
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("unmarshal gotify message failed: %w", err)
	}
	return &msg, nil
}

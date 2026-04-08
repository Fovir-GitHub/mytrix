package service

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Fovir-GitHub/mytrix/internal/config"
	"github.com/Fovir-GitHub/mytrix/internal/model"
)

type GotifyService struct{}

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

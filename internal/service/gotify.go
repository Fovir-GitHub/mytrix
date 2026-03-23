package service

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Fovir-GitHub/mytrix/internal/config"
	myhttp "github.com/Fovir-GitHub/mytrix/internal/http"
	"github.com/Fovir-GitHub/mytrix/internal/model"
)

type GotifyService interface {
	FetchMessages() ([]model.GotifyMessage, error)
}

type NoopGotifyService struct{}

type RealGotifyService struct {
	client *myhttp.Client
	server string
	token  string
}

func newGotifyService(client *myhttp.Client) GotifyService {
	cfg := config.Config.Gotify
	if !cfg.Enabled {
		slog.Info("gotify disabled")
		return &NoopGotifyService{}
	}

	if cfg.Server == "" || cfg.Token == "" {
		slog.Error(
			"gotify config invalid",
			"err", fmt.Errorf("server or token is empty"),
		)
		return &NoopGotifyService{}
	}

	slog.Info(
		"gotify service initialized",
		"server", cfg.Server,
	)
	return &RealGotifyService{
		client: client,
		server: cfg.Server,
		token:  cfg.Token,
	}
}

func (n *NoopGotifyService) FetchMessages() ([]model.GotifyMessage, error) {
	return []model.GotifyMessage{}, nil
}

func (gs *RealGotifyService) FetchMessages() ([]model.GotifyMessage, error) {
	start := time.Now()

	slog.Debug(
		"fetch gotify messages start",
		"server", gs.server,
	)

	req, err := http.NewRequest("GET", gs.server+"/message", nil)
	if err != nil {
		return nil, fmt.Errorf("create gotify request failed: %w", err)
	}
	req.Header.Set("X-Gotify-Key", gs.token)

	resp, err := gs.client.Do(req)
	if err != nil {
		slog.Debug(
			"gotify request failed",
			"server", gs.server,
			"error", err,
			"duration", time.Since(start),
		)
		return nil, fmt.Errorf("get gotify response failed: %w", err)
	}
	defer resp.Body.Close()

	var data struct {
		Messages []model.GotifyMessage `json:"messages"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode gotify json failed: %w", err)
	}

	slog.Debug(
		"gotify message fetched",
		"len", len(data.Messages),
		"duration", time.Since(start),
	)

	return data.Messages, nil
}

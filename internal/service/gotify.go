package service

import (
	"encoding/json"
	"fmt"
	"net/http"

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
		return &NoopGotifyService{}
	}

	return &RealGotifyService{
		client: client,
		server: cfg.Server,
		token:  cfg.Token,
	}
}

func (n *NoopGotifyService) FetchMessages() ([]model.GotifyMessage, error) {
	return nil, nil
}

func (gs *RealGotifyService) FetchMessages() ([]model.GotifyMessage, error) {
	req, err := http.NewRequest("GET", gs.server+"/message", nil)
	if err != nil {
		return nil, fmt.Errorf("create gotify request failed: %w", err)
	}
	req.Header.Set("X-Gotify-Key", gs.token)

	resp, err := gs.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get gotify response failed: %w", err)
	}
	defer resp.Body.Close()

	var data struct {
		Messages []model.GotifyMessage `json:"messages"`
	}
	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		return nil, fmt.Errorf("decode gotify json failed: %w", err)
	}
	return data.Messages, nil
}

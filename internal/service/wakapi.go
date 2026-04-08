package service

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/Fovir-GitHub/mytrix/internal/config"
	"github.com/Fovir-GitHub/mytrix/internal/http"
	"github.com/Fovir-GitHub/mytrix/internal/model"
	"github.com/Fovir-GitHub/mytrix/internal/scheduler"
)

type WakapiService interface {
	fetchData(model.WakapiInterval) (*model.WakapiData, error)
	FetchReport(model.WakapiInterval) (string, error)
}

type NoopWakapiService struct {
	err error
}

type RealWakapiService struct {
	c      *http.Client
	server string
	apiKey string
	userID string
	s      *scheduler.Scheduler
}

func newWakapiService(c *http.Client, s *scheduler.Scheduler) WakapiService {
	cfg := config.Config.Wakapi
	if !cfg.Enabled {
		return &NoopWakapiService{
			err: fmt.Errorf("wakapi is not enabled"),
		}
	}
	slog.Info("wakapi enabled")
	return &RealWakapiService{
		c:      c,
		server: cfg.Server,
		apiKey: cfg.APIKey,
		userID: cfg.UserID,
		s:      s,
	}
}

func (w *NoopWakapiService) fetchData(model.WakapiInterval) (*model.WakapiData, error) {
	return nil, w.err
}

func (w *NoopWakapiService) FetchReport(model.WakapiInterval) (string, error) {
	return "", w.err
}

func (w *RealWakapiService) fetchData(interval model.WakapiInterval) (*model.WakapiData, error) {
	var data struct {
		Data model.WakapiData `json:"data"`
	}

	const basePath = "/api/compat/wakatime/v1/"
	u := &url.URL{
		Scheme: "https",
		Host:   w.server,
	}
	u = u.JoinPath(basePath, "users", w.userID, "stats", string(interval))

	auth := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(w.apiKey)))
	req, err := http.NewRequest(
		http.MethodGet,
		u.String(),
		nil,
		map[string]string{
			"Authorization": auth,
		})
	if err != nil {
		return nil, fmt.Errorf("wakapi create http request failed: %w", err)
	}
	if err := w.c.DoJSON(req, &data); err != nil {
		return nil, fmt.Errorf("get json failed: %w", err)
	}
	return &data.Data, nil
}

func (w *RealWakapiService) FetchReport(interval model.WakapiInterval) (string, error) {
	data, err := w.fetchData(interval)
	if err != nil {
		return "", fmt.Errorf("fetch wakapi data failed: %w", err)
	}
	return data.ToMarkdown(), nil
}

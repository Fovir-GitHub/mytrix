package service

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/Fovir-GitHub/mytrix/internal/config"
	myhttp "github.com/Fovir-GitHub/mytrix/internal/http"
	"github.com/Fovir-GitHub/mytrix/internal/model"
	"github.com/Fovir-GitHub/mytrix/internal/scheduler"
)

type WakapiService interface {
	FetchData(model.WakapiInterval) (*model.WakapiData, error)
}

type NoopWakapiService struct{}

type RealWakapiService struct {
	c      *myhttp.Client
	server string
	apiKey string
	userID string
	s      *scheduler.Scheduler
}

func newWakapiService(c *myhttp.Client, s *scheduler.Scheduler) WakapiService {
	cfg := config.Config.Wakapi
	if !cfg.Enabled {
		return &NoopWakapiService{}
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

func (w *NoopWakapiService) FetchData(model.WakapiInterval) (*model.WakapiData, error) {
	return nil, fmt.Errorf("wakapi is not enabled")
}

func (w *RealWakapiService) FetchData(interval model.WakapiInterval) (*model.WakapiData, error) {
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
	req, err := myhttp.NewRequest(
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

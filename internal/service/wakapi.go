package service

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"sort"

	"github.com/Fovir-GitHub/mytrix/internal/config"
	myhttp "github.com/Fovir-GitHub/mytrix/internal/http"
	"github.com/Fovir-GitHub/mytrix/internal/model"
	"github.com/Fovir-GitHub/mytrix/internal/utils"
)

type WakapiService interface {
	FetchLanguages(model.WakapiInterval) ([]model.WakapiLanguage, error)
}

type NoopWakapiService struct{}

type RealWakapiService struct {
	c      *myhttp.Client
	server string
	apiKey string
	userID string
}

func newWakapiService(c *myhttp.Client) WakapiService {
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
	}
}

func (w *NoopWakapiService) FetchLanguages(model.WakapiInterval) ([]model.WakapiLanguage, error) {
	return nil, fmt.Errorf("wakapi is not enabled")
}

func (w *RealWakapiService) FetchLanguages(interval model.WakapiInterval) ([]model.WakapiLanguage, error) {
	var data struct {
		Data struct {
			Languages []model.WakapiLanguage `json:"languages"`
		} `json:"data"`
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

	langs := data.Data.Languages
	langs = utils.Filter(langs, func(w *model.WakapiLanguage) bool { return w.Percent >= 0.01 })
	sort.Slice(langs, func(i, j int) bool { return langs[i].Percent > langs[j].Percent })

	return langs, nil
}

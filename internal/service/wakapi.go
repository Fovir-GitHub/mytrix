package service

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/url"

	"codeberg.org/Fovir/mytrix/internal/config"
	"codeberg.org/Fovir/mytrix/internal/http"
	"codeberg.org/Fovir/mytrix/internal/model"
	"codeberg.org/Fovir/mytrix/internal/scheduler"
	"codeberg.org/Fovir/mytrix/internal/utils"
)

// WakapiService interface defines methods for Wakapi service implementations.
// It provides methods to fetch Wakapi data and generate reports.
type WakapiService interface {
	// fetchData retrieves Wakapi data for the given interval.
	// It returns the WakapiData and any error encountered.
	fetchData(model.WakapiInterval) (*model.WakapiData, error)

	// FetchReport generates a formatted report for the given Wakapi interval.
	// It returns the report as a string and any error encountered.
	FetchReport(model.WakapiInterval) (string, error)
}

// NoopWakapiService is a WakapiService implementation that returns errors.
// It is used when Wakapi is disabled in the configuration.
type NoopWakapiService struct {
	err error
}

// RealWakapiService implements WakapiService for real Wakapi API interactions.
// It holds the HTTP client, server configuration, API key, user ID, and scheduler.
type RealWakapiService struct {
	c      *http.Client
	server string
	apiKey string
	userID string
	s      *scheduler.Scheduler
}

// NewWakapiService creates a WakapiService based on configuration.
// If Wakapi is enabled, it returns a RealWakapiService; otherwise, it returns a NoopWakapiService.
// It takes an HTTP client and scheduler as dependencies.
func NewWakapiService(c *http.Client, s *scheduler.Scheduler) WakapiService {
	cfg := config.Config.Wakapi
	if !cfg.Enabled {
		slog.Info("wakapi disabled")
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

// fetchData returns nil and the stored error for any interval.
// It is used when Wakapi is disabled to simulate a service failure.
func (w *NoopWakapiService) fetchData(model.WakapiInterval) (*model.WakapiData, error) {
	return nil, w.err
}

// FetchReport returns an empty string and the stored error for any interval.
// It is used when Wakapi is disabled to simulate a service failure.
func (w *NoopWakapiService) FetchReport(model.WakapiInterval) (string, error) {
	return "", w.err
}

// fetchData retrieves Wakapi data for the specified interval from the Wakapi API.
// It constructs the request URL, adds authentication, and processes the response.
// Language data is filtered to only include entries with percentage >= 0.01%.
// Returns the WakapiData and any error encountered during the process.
func (w *RealWakapiService) fetchData(interval model.WakapiInterval) (*model.WakapiData, error) {
	slog.Debug("fetch wakapi data start", "interval", string(interval))
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
		return nil, fmt.Errorf("fetch wakapi data failed: %w", err)
	}
	if err := w.c.DoJSON(req, &data); err != nil {
		return nil, fmt.Errorf("fetch wakapi data failed: %w", err)
	}
	slog.Debug("wakapi data fetched")

	data.Data.Langs = utils.Filter(data.Data.Langs, func(lang *model.WakapiLanguage) bool {
		return lang.Percent >= 0.01
	})

	return &data.Data, nil
}

// FetchReport generates a formatted report for the given Wakapi interval.
// It fetches the Wakapi data and converts it to markdown format.
// Returns the formatted report string and any error encountered.
func (w *RealWakapiService) FetchReport(interval model.WakapiInterval) (string, error) {
	slog.Debug("fetch wakapi report start", "interval", string(interval))
	data, err := w.fetchData(interval)
	if err != nil {
		return "", fmt.Errorf("fetch wakapi report failed (interval=%s): %w", string(interval), err)
	}
	return data.ToMarkdown(), nil
}

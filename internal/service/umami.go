package service

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"codeberg.org/Fovir/mytrix/internal/config"
	myhttp "codeberg.org/Fovir/mytrix/internal/http"
	"codeberg.org/Fovir/mytrix/internal/model"
)

// UmamiService interface defines methods for Umami service implementations.
// It provides methods for authentication, fetching website data, statistics, and generating reports.
type UmamiService interface {
	FetchReport(*model.UmamiInterval) string
}

type RealUmamiService struct {
	c        *myhttp.Client
	token    string
	server   string
	username string
	password string
}

func NewUmamiService(c *myhttp.Client) UmamiService {
	cfg := config.Config.Umami
	slog.Info("umami initialized", "enabled", cfg.Enabled, "server", cfg.Server)

	noop := &NoopUmamiService{err: fmt.Errorf("umami is not enabled")}
	if !cfg.Enabled {
		return noop
	}

	us := &RealUmamiService{
		c:        c,
		server:   cfg.Server,
		username: cfg.Username,
		password: cfg.Password,
	}
	t, err := us.getToken()
	if err != nil {
		slog.Warn("get umami token failed, use empty token", "err", err)
		t = ""
	}
	us.token = t
	return us
}

func (ru *RealUmamiService) fetchWebsites() ([]*model.UmamiWebsite, error) {
	slog.Debug("fetch umami websites start")

	var data struct {
		Data []*model.UmamiWebsite `json:"data"`
	}

	u := ru.createURL("/api/websites")
	req, err := myhttp.NewRequest(
		myhttp.MethodGet,
		u.String(),
		nil,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("fetch umami websites failed: %w", err)
	}
	ru.setAuthHeader(req)

	err = ru.c.DoJSON(req, &data)
	if err != nil {
		slog.Warn("umami fetch websites failed, retry with new token", "err", err)
		ru.updateToken()
		if err := ru.c.DoJSON(req, &data); err != nil {
			return nil, fmt.Errorf("fetch umami websites failed: %w", err)
		}
	}
	websites := data.Data
	slog.Debug("got umami websites", "len", len(websites))
	return websites, nil
}

func (ru *RealUmamiService) fetchWebsiteStat(website *model.UmamiWebsite, interval *model.UmamiInterval) (*model.UmamiWebsiteStat, error) {
	slog.Debug("fetch umami website stat begin", "name", website.Name)

	var stat *model.UmamiWebsiteStat
	const basePath = "/api/websites"
	u := ru.createURL("")
	u = u.JoinPath(basePath, website.ID, "stats")
	q := u.Query()
	q.Set("startAt", fmt.Sprintf("%d", interval.Start.UnixMilli()))
	q.Set("endAt", fmt.Sprintf("%d", interval.End.UnixMilli()))
	u.RawQuery = q.Encode()

	req, err := myhttp.NewRequest(myhttp.MethodGet, u.String(), nil, nil)
	if err != nil {
		return nil, fmt.Errorf("fetch umami website stat failed (name=%s, domain=%s): %w", website.Name, website.Domain, err)
	}
	ru.setAuthHeader(req)

	if err := ru.c.DoJSON(req, &stat); err != nil {
		return nil, fmt.Errorf("fetch umami website stat failed (name=%s, domain=%s): %w", website.Name, website.Domain, err)
	}
	slog.Debug("fetched umami website stat", "name", website.Name, "visits", stat.Visitors)
	return stat, nil
}

func (ru *RealUmamiService) fetchWebsiteData(interval *model.UmamiInterval) ([]*model.UmamiWebsite, error) {
	websites, err := ru.fetchWebsites()
	if err != nil {
		return nil, fmt.Errorf("fetch umami website data failed: %w", err)
	}

	var (
		res  []*model.UmamiWebsite
		errs []error
	)
	for _, w := range websites {
		stat, err := ru.fetchWebsiteStat(w, interval)
		if err != nil {
			slog.Warn("fetch website data failed", "name", w.Name, "err", err)
			errs = append(errs, fmt.Errorf("fetch website data failed: %w", err))
			continue
		}
		w.Stat = stat
		res = append(res, w)
	}

	slog.Info("umami website data fetched", "len", len(res))
	if len(errs) > 0 {
		joined := errors.Join(errs...)
		slog.Warn("fetch umami website data partially failed",
			"err", joined)
		return res, joined
	}

	return res, nil
}

func (ru *RealUmamiService) FetchReport(interval *model.UmamiInterval) string {
	slog.Debug("fetch umami report start", "interval", interval)
	websites, err := ru.fetchWebsiteData(interval)
	if err != nil {
		return err.Error()
	}
	return ru.generateReport(websites)
}

func (ru RealUmamiService) generateReport(websites []*model.UmamiWebsite) string {
	var report strings.Builder
	for _, w := range websites {
		report.WriteString(w.ToMarkdown())
		report.WriteString("\n")
	}
	return report.String()
}

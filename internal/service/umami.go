package service

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/Fovir-GitHub/mytrix/internal/config"
	myhttp "github.com/Fovir-GitHub/mytrix/internal/http"
	"github.com/Fovir-GitHub/mytrix/internal/model"
)

type UmamiService interface {
	getToken() (string, error)
	fetchWebsites() ([]*model.UmamiWebsite, error)
	fetchWebsiteStat(*model.UmamiWebsite, *model.UmamiInterval) (*model.UmamiWebsiteStat, error)
	fetchWebsiteData(*model.UmamiInterval) ([]*model.UmamiWebsite, error)
	generateReport([]*model.UmamiWebsite) string
	FetchReport(*model.UmamiInterval) string
}

type RealUmamiService struct {
	c        *myhttp.Client
	token    string
	server   string
	username string
	password string
}

func newUmamiService(c *myhttp.Client) UmamiService {
	cfg := config.Config.Umami
	noop := &NoopUmamiService{err: fmt.Errorf("umami is not enabled")}
	if !cfg.Enabled {
		slog.Info("umami disabled")
		return noop
	}

	slog.Info("umami enabled")
	us := &RealUmamiService{
		c:        c,
		server:   cfg.Server,
		username: cfg.Username,
		password: cfg.Password,
	}
	t, err := us.getToken()
	if err != nil {
		slog.Error("get token failed, umami is disabled", "err", err)
		return noop
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
		return nil, fmt.Errorf("create umami fetch websites request failed: %w", err)
	}
	ru.setAuthHeader(req)

	err = ru.c.DoJSON(req, &data)
	if err != nil {
		slog.Warn("umami fetch websites failed, retry to login", "err", err)
		ru.updateToken()
		if err := ru.c.DoJSON(req, &data); err != nil {
			return nil, fmt.Errorf("umami fetch websites failed: %w", err)
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
		return nil, fmt.Errorf("create umami fetch website stat request failed: %w", err)
	}
	ru.setAuthHeader(req)

	if err := ru.c.DoJSON(req, &stat); err != nil {
		return nil, fmt.Errorf("umami fetch website stat failed: %w", err)
	}
	slog.Debug("fetched umami website stat", "stat", stat)
	return stat, nil
}

func (ru *RealUmamiService) fetchWebsiteData(interval *model.UmamiInterval) ([]*model.UmamiWebsite, error) {
	slog.Debug("fetch umami website data begin")

	websites, err := ru.fetchWebsites()
	if err != nil {
		return nil, fmt.Errorf("fetch umami websites failed: %w", err)
	}

	var res []*model.UmamiWebsite
	for _, w := range websites {
		stat, err := ru.fetchWebsiteStat(w, interval)
		if err != nil {
			return nil, fmt.Errorf("fetch website stat failed: %w", err)
		}
		w.Stat = stat
		res = append(res, w)
	}
	slog.Debug("umami website data fetched", "len", len(res))
	return res, nil
}

func (ru *RealUmamiService) FetchReport(interval *model.UmamiInterval) string {
	slog.Debug("fetch umami report start", "interval", interval)
	websites, err := ru.fetchWebsiteData(interval)
	if err != nil {
		slog.Error("fetch umami website data failed", "err", err)
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

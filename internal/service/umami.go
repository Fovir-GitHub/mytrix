package service

import (
	"fmt"
	"log/slog"

	"github.com/Fovir-GitHub/mytrix/internal/config"
	myhttp "github.com/Fovir-GitHub/mytrix/internal/http"
	"github.com/Fovir-GitHub/mytrix/internal/model"
)

type UmamiService interface {
	getToken() (string, error)
	FetchWebsites() ([]model.UmamiWebsite, error)
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
	if !cfg.Enabled {
		return &NoopUmamiService{err: fmt.Errorf("umami is not enabled")}
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
		slog.Error("get token failed", "err", err)
		return &NoopUmamiService{}
	}
	us.token = t
	return us
}

func (ru *RealUmamiService) FetchWebsites() ([]model.UmamiWebsite, error) {
	var data struct {
		Data []model.UmamiWebsite `json:"data"`
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
	return data.Data, nil
}

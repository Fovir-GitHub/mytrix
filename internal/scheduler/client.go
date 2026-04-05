package scheduler

import (
	"log/slog"
	"time"
	_ "time/tzdata"

	"github.com/Fovir-GitHub/mytrix/internal/config"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	c *cron.Cron
}

func NewScheduler() *Scheduler {
	slog.Info("created scheduler")

	cfgTz := config.Config.TZ
	tzOpt := cron.WithLocation(time.Local)
	if cfgTz != "" {
		loc, err := time.LoadLocation(cfgTz)
		if err != nil {
			slog.Warn(
				"invalid timezone, use default location",
				"tz", cfgTz,
				"default", time.Local.String(),
				"err", err,
			)
		} else {
			tzOpt = cron.WithLocation(loc)
		}
	}

	s := &Scheduler{
		c: cron.New(tzOpt),
	}
	slog.Info("set timezone", "timezone", s.c.Location().String())
	return s
}

func (s *Scheduler) Start() {
	s.c.Start()
}

func (s *Scheduler) Register(t string, job func()) {
	_, err := s.c.AddFunc(t, job)
	if err != nil {
		slog.Error("register schedule failed", "time", t, "err", err)
	}
}

package scheduler

import (
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	c *cron.Cron
}

func NewScheduler() *Scheduler {
	slog.Info("created scheduler")

	s := &Scheduler{
		c: cron.New(cron.WithLocation(time.Local)),
	}

	slog.Info("set scheduler timezone", "timezone", s.c.Location().String())
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

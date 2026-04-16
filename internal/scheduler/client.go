// Package scheduler provides job scheduling functionality.
// It uses the robfig/cron library to schedule and execute jobs.
package scheduler

import (
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"
)

// Scheduler manages the scheduling and execution of background jobs using cron expressions.
type Scheduler struct {
	c *cron.Cron
}

// NewScheduler creates a new Scheduler with local timezone configuration.
func NewScheduler() *Scheduler {
	slog.Info("created scheduler")

	s := &Scheduler{
		c: cron.New(cron.WithLocation(time.Local)),
	}

	slog.Info("set scheduler timezone", "timezone", s.c.Location().String())
	return s
}

// Start begins executing all registered scheduled jobs.
func (s *Scheduler) Start() {
	slog.Info("start the scheduler")
	s.c.Start()
}

// Register schedules a job to run at the time specified by the cron expression.
// The job function is called whenever the cron expression matches the current time.
func (s *Scheduler) Register(t string, job func()) {
	_, err := s.c.AddFunc(t, job)
	if err != nil {
		slog.Error("register schedule failed", "time", t, "err", err)
		return
	}
	slog.Debug("registered scheduler", "time", t)
}

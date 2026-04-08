package bot

import (
	"log/slog"

	"github.com/Fovir-GitHub/mytrix/internal/scheduler"
)

func (b *Bot) registerScheduler() {
	jobList := []func() []scheduler.ScheduledJob{
		b.Handler.WakapiScheduleList,
		b.Handler.UmamiScheduleList,
	}

	var jobs []scheduler.ScheduledJob
	for _, l := range jobList {
		jobs = append(jobs, l()...)
	}
	for _, j := range jobs {
		b.Scheduler.Register(j.Cron, j.Job)
	}
	slog.Info("schedulers registered")
}

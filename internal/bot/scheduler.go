package bot

import (
	"log/slog"

	"codeberg.org/Fovir/mytrix/internal/scheduler"
)

// registerScheduler registers schedulers from various modules by collecting scheduled jobs from
// Wakapi and Umami handlers and registering them with the scheduler.
// If a module has schedulers, it should provide a scheduler list and add it to jobList.
func (b *Bot) registerScheduler() {
	jobList := []func() []scheduler.ScheduledJob{
		b.Handler.WakapiScheduleList,
		b.Handler.UmamiScheduleList,
		b.Handler.RSSScheduleList,
	}

	var jobs []scheduler.ScheduledJob
	for _, l := range jobList {
		jobs = append(jobs, l()...)
	}
	for _, j := range jobs {
		b.Scheduler.Register(j.Cron, j.Job)
	}
	slog.Info("schedulers registered",
		"count", len(jobs))
}

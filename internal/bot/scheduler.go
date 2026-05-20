package bot

import (
	"errors"
	"log/slog"

	"codeberg.org/Fovir/mytrix/internal/scheduler"
)

// registerScheduler registers schedulers from various modules by collecting scheduled jobs from
// Wakapi, Umami and RSS handlers and registering them with the scheduler.
// If a module has schedulers, it should provide a scheduler list and add it to jobList.
func (b *Bot) registerScheduler() error {
	jobList := []func() []scheduler.ScheduledJob{
		b.Handler.WakapiScheduleList,
		b.Handler.UmamiScheduleList,
		b.Handler.RSSScheduleList,
	}

	var errs []error
	jobCounter := 0
	for _, l := range jobList {
		jobs := l()

		for _, j := range jobs {
			if err := b.Scheduler.Register(j.Cron, j.Job); err != nil {
				errs = append(errs, err)
				continue
			}
			jobCounter++
		}
	}
	slog.Info("schedulers registered",
		"count", jobCounter)

	return errors.Join(errs...)
}

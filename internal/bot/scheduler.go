package bot

import (
	"log/slog"

	"github.com/Fovir-GitHub/mytrix/internal/scheduler"
)

func (b *Bot) registerScheduler() {
	var jobs []scheduler.ScheduledJob
	jobs = append(jobs, b.Handler.WakapiScheduleList()...)
	for _, j := range jobs {
		b.Scheduler.Register(j.Cron, j.Job)
	}
	slog.Info("schedulers registered")
}

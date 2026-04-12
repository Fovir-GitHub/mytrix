package handler

import (
	"context"
	"log/slog"

	"github.com/Fovir-GitHub/mytrix/internal/config"
	"github.com/Fovir-GitHub/mytrix/internal/model"
	"github.com/Fovir-GitHub/mytrix/internal/scheduler"
	"maunium.net/go/mautrix/id"
)

// handleWakapiSchedule handles scheduled Wakapi report generation and sending.
// It fetches a Wakapi report for the specified interval and sends it to the configured Matrix room.
func (h *Handler) handleWakapiSchedule(ctx context.Context, interval model.WakapiInterval) {
	roomID := config.Config.RoomID
	report, err := h.service.Wakapi.FetchReport(interval)
	if err != nil {
		slog.Error("fetch wakapi report failed", "interval", interval, "err", err)
		return
	}
	_ = h.service.Message.Reply(ctx, id.RoomID(roomID), report)
}

// WakapiScheduleList returns a list of scheduled jobs for Wakapi report generation.
// It maps configured cron expressions to their corresponding Wakapi intervals:
// daily reports use yesterday's data, weekly reports use the last 7 days,
// monthly reports use the last 30 days, and yearly reports use the last 12 months.
// Each scheduled job invokes handleWakapiSchedule with its respective interval.
func (h *Handler) WakapiScheduleList() []scheduler.ScheduledJob {
	cfg := config.Config.Wakapi
	m := map[string]model.WakapiInterval{
		cfg.DailyReportCron:   model.WakapiIntervalYesterday,
		cfg.WeeklyReportCron:  model.WakapiIntervalLast7Days,
		cfg.MonthlyReportCron: model.WakapiIntervalLast30Days,
		cfg.YearlyReportCron:  model.WakapiIntervalLast12Months,
	}
	var res []scheduler.ScheduledJob
	for cron, interval := range m {
		sj := scheduler.ScheduledJob{
			Cron: cron,
			Job:  func() { h.handleWakapiSchedule(context.Background(), interval) },
		}
		res = append(res, sj)
	}
	return res
}

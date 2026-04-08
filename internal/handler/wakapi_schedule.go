package handler

import (
	"context"
	"log/slog"

	"github.com/Fovir-GitHub/mytrix/internal/config"
	"github.com/Fovir-GitHub/mytrix/internal/model"
	"github.com/Fovir-GitHub/mytrix/internal/scheduler"
	"maunium.net/go/mautrix/id"
)

func (h *Handler) handleWakapiSchedule(ctx context.Context, interval model.WakapiInterval) {
	roomID := config.Config.RoomID
	report, err := h.service.Wakapi.FetchReport(interval)
	if err != nil {
		slog.Error("fetch wakapi report failed", "interval", interval, "err", err)
		return
	}
	_ = h.service.Message.Reply(ctx, id.RoomID(roomID), report)
}

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

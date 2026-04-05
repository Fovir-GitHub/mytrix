package handler

import (
	"context"
	"log/slog"

	"github.com/Fovir-GitHub/mytrix/internal/config"
	"github.com/Fovir-GitHub/mytrix/internal/model"
	"github.com/Fovir-GitHub/mytrix/internal/scheduler"
	"maunium.net/go/mautrix/id"
)

func (h *Handler) handleWakapiSchedule(ctx context.Context, interval model.WakapiInterval, text string) {
	roomID := config.Config.RoomID
	report, err := h.fetchWakapiReport(interval)
	if err != nil {
		slog.Error("fetch wakapi report failed", "err", err)
		return
	}
	msg := text + "\n" + report
	if err := h.service.Message.Reply(ctx, id.RoomID(roomID), msg); err != nil {
		slog.Error(
			"send message failed",
			"report", report,
			"roomID", roomID,
			"err", err,
		)
	}
}

func (h *Handler) WakapiScheduleList() []scheduler.ScheduledJob {
	cfg := config.Config.Wakapi
	return []scheduler.ScheduledJob{
		{
			Cron: cfg.DailyReportCron,
			Job:  func() { h.handleWakapiSchedule(context.Background(), model.WakapiIntervalYesterday, "Daily Report") },
		},
		{
			Cron: cfg.MonthlyReportCron,
			Job:  func() { h.handleWakapiSchedule(context.Background(), model.WakapiIntervalLast30Days, "Monthly Report") },
		},
		{
			Cron: cfg.YearlyReportCron,
			Job: func() {
				h.handleWakapiSchedule(context.Background(), model.WakapiIntervalLast12Months, "Yearly Report")
			},
		},
	}
}

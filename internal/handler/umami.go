package handler

import (
	"context"
	"log/slog"
	"strings"

	"github.com/Fovir-GitHub/mytrix/internal/config"
	"github.com/Fovir-GitHub/mytrix/internal/model"
	"github.com/Fovir-GitHub/mytrix/internal/scheduler"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

func (h *Handler) handleUmamiCommand(ctx context.Context, evt *event.Event) error {
	slog.Debug("handle umami command start")
	interval := getUmamiInterval(evt.Content.AsMessage().Body)
	report := h.service.Umami.FetchReport(interval)
	return h.service.Message.Reply(ctx, evt.RoomID, report)
}

func getUmamiInterval(msg string) *model.UmamiInterval {
	defaultInterval := model.UmamiIntervalMap[config.Config.Umami.DefaultInterval]()
	slog.Debug("start getUmamiInterval", "defaultInterval", config.Config.Umami.DefaultInterval)
	parts := strings.Fields(msg)
	if len(parts) < 2 {
		slog.Warn("no interval provided, fallback to default interval", "defaultInterval", defaultInterval)
		return defaultInterval
	}
	intervalStr := parts[1]
	slog.Debug("get umami interval", "intervalStr", intervalStr)
	interval, err := model.ParseUmamiInterval(intervalStr)
	if err != nil {
		slog.Warn(
			"parse umami interval failed, use default interval",
			"msg", msg,
			"intervalStr", intervalStr,
			"defaultInterval", defaultInterval,
			"err", err,
		)
		return defaultInterval
	}
	slog.Debug("got umami interval", "interval", interval)
	return interval
}

func (h *Handler) handleUmamiSchedule(ctx context.Context, interval *model.UmamiInterval) {
	slog.Debug("handle umami schedule", "interval", interval)
	roomID := config.Config.RoomID
	report := h.service.Umami.FetchReport(interval)
	_ = h.service.Message.Reply(ctx, id.RoomID(roomID), report)
}

func (h *Handler) UmamiScheduleList() []scheduler.ScheduledJob {
	cfg := config.Config.Umami
	m := map[string]func() *model.UmamiInterval{
		cfg.DailyReportCron:   model.UmamiIntervalYesterday,
		cfg.WeeklyReportCron:  model.UmamiIntervalLastWeek,
		cfg.MonthlyReportCron: model.UmamiIntervalLastMonth,
		cfg.YearlyReportCron:  model.UmamiIntervalLastYear,
	}

	var res []scheduler.ScheduledJob
	for cron, interval := range m {
		res = append(res, scheduler.ScheduledJob{
			Cron: cron,
			Job:  func() { h.handleUmamiSchedule(context.Background(), interval()) },
		})
	}
	return res
}

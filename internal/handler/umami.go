package handler

import (
	"context"
	"log/slog"
	"strings"

	"codeberg.org/Fovir/mytrix/internal/config"
	"codeberg.org/Fovir/mytrix/internal/model"
	"codeberg.org/Fovir/mytrix/internal/scheduler"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

// handleUmamiCommand processes the !umami command from a Matrix event.
// It extracts the interval from the message content, fetches the corresponding Umami report, and sends it to the room where the command was issued.
func (h *Handler) handleUmamiCommand(ctx context.Context, evt *event.Event) error {
	interval := getUmamiInterval(evt.Content.AsMessage().Body)
	slog.Debug("handle umami command start",
		"start", interval.Start.String(),
		"end", interval.End.String())
	report := h.service.Umami.FetchReport(interval)
	return h.service.Message.Reply(ctx, evt.RoomID, report)
}

// getUmamiInterval extracts the Umami interval from the given message string.
// It splits the message into fields and checks the second field for a valid interval.
// If no interval is provided or the provided interval is invalid, it falls back to
// the configured default interval. It returns a pointer to the selected UmamiInterval.
func getUmamiInterval(msg string) *model.UmamiInterval {
	defaultInterval := model.UmamiIntervalMap[config.Config.Umami.DefaultInterval]()
	parts := strings.Fields(msg)
	if len(parts) < 2 {
		slog.Warn("no interval provided, fallback to default interval", "defaultInterval", defaultInterval)
		return defaultInterval
	}
	intervalStr := parts[1]
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
	slog.Debug("got umami interval", "input", intervalStr, "resolved", interval)
	return interval
}

// handleUmamiSchedule handles scheduled Umami report generation and sending.
// It fetches an Umami report for the specified interval and sends it to the configured Matrix room.
func (h *Handler) handleUmamiSchedule(ctx context.Context, interval *model.UmamiInterval) {
	slog.Debug("handle umami schedule", "interval", interval)
	roomID := config.Config.RoomID
	report := h.service.Umami.FetchReport(interval)
	_ = h.service.Message.Reply(ctx, id.RoomID(roomID), report)
}

// UmamiScheduleList returns a list of scheduled jobs for Umami report generation.
// It maps configured cron expressions to their corresponding interval functions,
// creating scheduled jobs for daily, weekly, monthly, and yearly reports.
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

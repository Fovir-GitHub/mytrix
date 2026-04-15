package handler

import (
	"context"
	"log/slog"
	"strings"

	"codeberg.org/Fovir/mytrix/internal/config"
	"codeberg.org/Fovir/mytrix/internal/scheduler"
	"maunium.net/go/mautrix/id"
)

func (h *Handler) handleRSSSchedule(ctx context.Context) {
	roomID := config.Config.RoomID
	updated, err := h.service.RSS.Update()
	if err != nil {
		slog.Debug("update rss error", "err", err)
	}
	if len(updated) <= 0 {
		return
	}

	var msg strings.Builder
	for _, item := range updated {
		msg.WriteString(item.ToMarkdown())
		msg.WriteString("\n")
	}
	_ = h.service.Message.Reply(ctx, id.RoomID(roomID), msg.String())
}

func (h *Handler) RSSScheduleList() []scheduler.ScheduledJob {
	cfg := config.Config.RSS
	return []scheduler.ScheduledJob{
		{
			Cron: cfg.Cron,
			Job: func() {
				h.handleRSSSchedule(context.Background())
			},
		},
	}
}

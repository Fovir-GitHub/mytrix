package handler

import (
	"context"
	"log/slog"

	"codeberg.org/Fovir/mytrix/internal/config"
	"codeberg.org/Fovir/mytrix/internal/scheduler"
	"maunium.net/go/mautrix/id"
)

func (h *Handler) handleRSSSchedule(ctx context.Context) {
	roomID := config.Config.RoomID
	updated, err := h.service.RSS.Update()
	if err != nil {
		slog.Error("update rss error",
			"room_id", roomID, "err", err)
	}

	slog.Info("rss schedule update done", "items", len(updated))
	if len(updated) <= 0 {
		return
	}

	_ = h.service.Message.Reply(ctx, id.RoomID(roomID), updated)
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

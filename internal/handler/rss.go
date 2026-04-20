package handler

import (
	"context"
	"errors"
	"log/slog"

	"codeberg.org/Fovir/mytrix/internal/config"
	"codeberg.org/Fovir/mytrix/internal/scheduler"
	"codeberg.org/Fovir/mytrix/internal/service"
	"maunium.net/go/mautrix/id"
)

func (h *Handler) handleRSSSchedule(ctx context.Context) {
	roomID := config.Config.RoomID
	updated, err := h.service.RSS.Update()
	if err != nil {
		if errors.Is(err, service.ErrRSSFetchFeeds) {
			slog.Error("update rss error",
				"room_id", roomID, "err", err)
			return
		}

		if errors.Is(err, service.ErrRSSNoUpdate) {
			slog.Info("rss everything up to date")
			return
		}
	}

	slog.Info("rss schedule update done", "items", len(updated))

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

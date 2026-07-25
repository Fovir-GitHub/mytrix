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
	updated, err := h.service.RSS.Update(ctx)
	if err != nil {
		if errors.Is(err, service.ErrRSSFetchFeeds) {
			slog.Error("update rss error",
				"room_id", roomID, "err", err)
			return
		}

		if errors.Is(err, service.ErrRSSNoUpdate) {
			slog.Debug("rss everything up to date")
			return
		}

		if errors.Is(err, service.ErrRSSPartialUpdate) {
			slog.Warn("rss feed partially updated", "err", err)
		}
	}

	slog.Info("rss schedule update done", "items", len(updated))

	for _, item := range updated {
		if err := h.service.Message.ReplyWithoutResp(ctx, id.RoomID(roomID), item.Rendered); err != nil {
			slog.Error("send updated rss items failed",
				"room_id", roomID,
				"len", len(item.ItemIDs),
				"err", err)
			continue
		}

		if err := h.service.RSS.MarkItemsPushed(ctx, &item); err != nil {
			slog.Error("mark rss item pushed failed",
				"room_id", roomID,
				"len", len(item.ItemIDs),
				"err", err)
		}
	}
}

func (h *Handler) RSSScheduleList() []scheduler.ScheduledJob {
	cfg := config.Config.RSS
	if !cfg.Enabled {
		return nil
	}
	return []scheduler.ScheduledJob{
		{
			Cron: cfg.Cron,
			Job: func() {
				h.handleRSSSchedule(context.Background())
			},
		},
	}
}

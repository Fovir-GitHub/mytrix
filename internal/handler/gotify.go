package handler

import (
	"context"
	"log/slog"

	"codeberg.org/Fovir/mytrix/internal/config"
	"codeberg.org/Fovir/mytrix/internal/model"
	"maunium.net/go/mautrix/id"
)

func (h *Handler) handleGotify(ctx context.Context, event *model.WsEvent) error {
	msg, err := h.service.Gotify.HandleEvent(event)
	if err != nil {
		slog.Error("handle gotify event failed",
			"err", err)
		return err
	}
	return h.service.Message.Reply(ctx, id.RoomID(config.Config.RoomID), msg.ToMarkdown())
}

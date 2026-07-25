package handler

import (
	"context"
	"log/slog"

	"codeberg.org/Fovir/mytrix/internal/config"
	"codeberg.org/Fovir/mytrix/internal/model"
	"codeberg.org/Fovir/mytrix/internal/render"
	"maunium.net/go/mautrix/id"
)

// handleGotify processes incoming Gotify WebSocket events and sends the notification as a Matrix message.
func (h *Handler) handleGotify(ctx context.Context, event *model.WsEvent) error {
	msg, err := h.service.Gotify.HandleEvent(event)
	if err != nil {
		slog.Error("handle gotify event failed",
			"err", err)
		return err
	}

	return h.service.Message.ReplyWithoutResp(ctx, id.RoomID(config.Config.RoomID), render.GotifyMarkdown(msg))
}

package handler

import (
	"context"

	"codeberg.org/Fovir/mytrix/internal/config"
	"codeberg.org/Fovir/mytrix/internal/model"
	"maunium.net/go/mautrix/id"
)

func (h *Handler) handleGotify(ctx context.Context, event *model.WsEvent) error {
	msg, err := h.service.Gotify.HandleEvent(event)
	if err != nil {
		return err
	}
	return h.service.Message.Reply(ctx, id.RoomID(config.Config.RoomID), msg.ToMarkdown())
}

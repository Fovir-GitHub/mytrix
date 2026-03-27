package handler

import (
	"context"
	"fmt"

	"github.com/Fovir-GitHub/mytrix/internal/config"
	"github.com/Fovir-GitHub/mytrix/internal/model"
	"maunium.net/go/mautrix/id"
)

func (h *Handler) handleGotify(ctx context.Context, event *model.WsEvent) error {
	msg, err := h.service.Gotify.HandleEvent(event)
	if err != nil {
		return fmt.Errorf("handle gotify event failed: %w", err)
	}
	return h.service.Message.Reply(ctx, id.RoomID(config.Config.RoomID), msg.Message)
}

package handler

import (
	"context"
	"log/slog"

	"maunium.net/go/mautrix/event"
)

func (h *Handler) handlePing(ctx context.Context, evt *event.Event) error {
	slog.Debug("handle ping command")
	if err := h.service.Message.Ping(ctx, evt); err != nil {
		slog.Error("handle ping failed",
			"room", evt.RoomID.String(),
			"sender", evt.Sender.String(),
			"err", err)
		return err
	}
	return nil
}

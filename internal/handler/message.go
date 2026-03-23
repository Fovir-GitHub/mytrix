package handler

import (
	"context"
	"log/slog"

	"maunium.net/go/mautrix/event"
)

func (h *Handler) handlePing(ctx context.Context, evt *event.Event) error {
	slog.Debug("handle ping command")
	return h.service.Message.Ping(ctx, evt)
}

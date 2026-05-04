package handler

import (
	"context"
	"log/slog"

	"maunium.net/go/mautrix/event"
)

// handlePing responds to a ping command with "pong".
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

// handleVersion responds to "!version" command with current version.
func (h *Handler) handleVersion(ctx context.Context, evt *event.Event) error {
	slog.Debug("handle version command")
	err := h.service.Message.Version(ctx, evt)
	if err != nil {
		slog.Error("handle version failed",
			"room", evt.RoomID.String(),
			"sender", evt.Sender.String(),
			"err", err)
	}
	return err
}

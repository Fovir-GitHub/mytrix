package handler

import (
	"context"
	"log/slog"

	"maunium.net/go/mautrix/event"
)

func (h *Handler) HandleInvitation(ctx context.Context, evt *event.Event) {
	if err := h.service.Invitation.Process(ctx, evt); err != nil {
		slog.Error("handle invitation failed", "err", err)
		return
	}
	slog.Info("joined room", "room_id", evt.RoomID.String(), "sender", evt.Sender.String())
}

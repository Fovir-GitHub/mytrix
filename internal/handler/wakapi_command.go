// Package handler handles incoming events and commands.
package handler

import (
	"context"
	"errors"
	"log/slog"

	"maunium.net/go/mautrix/event"
)

// handleWakapiCommand processes the !wakapi command from a Matrix event.
// It extracts the time interval from the message content, fetches the corresponding Wakapi report, and sends it to the room where the command was issued.
// If the interval cannot be parsed, it sends the error message back to the room.
func (h *Handler) handleWakapiCommand(ctx context.Context, evt *event.Event) error {
	interval, err := getWakapiInterval(evt.Content.AsMessage().Body)
	if err != nil {
		slog.Error("get wakapi interval failed, reply errors", "err", err)
		replyErr := h.service.Message.Reply(ctx, evt.RoomID, err.Error())
		return errors.Join(err, replyErr)
	}

	report, err := h.service.Wakapi.FetchReport(interval)
	if err != nil {
		slog.Error("fetch wakapi report failed", "err", err)
		return err
	}
	return h.service.Message.Reply(ctx, evt.RoomID, report)
}

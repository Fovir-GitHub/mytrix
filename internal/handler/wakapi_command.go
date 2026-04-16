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
	reply := h.getReply(ctx, evt)
	if err != nil {
		slog.Error("get wakapi interval failed, reply errors", "err", err)
		replyErr := reply("Invalid Wakapi interval")
		return errors.Join(err, replyErr)
	}

	report, err := h.service.Wakapi.FetchReport(interval)
	if err != nil {
		slog.Error("fetch wakapi report failed", "err", err)
		return reply("Failed to fetch Wakapi report")
	}
	return reply(report)
}

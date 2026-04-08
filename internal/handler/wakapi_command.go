package handler

import (
	"context"
	"errors"
	"log/slog"

	"maunium.net/go/mautrix/event"
)

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

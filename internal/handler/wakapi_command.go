package handler

import (
	"context"
	"log/slog"

	"maunium.net/go/mautrix/event"
)

func (h *Handler) handleWakapiCommand(ctx context.Context, evt *event.Event) error {
	interval, err := getWakapiInterval(evt.Content.AsMessage().Body)
	if err != nil {
		slog.Error("get wakapi interval failed, reply errors", "err", err)
		if replyErr := h.service.Message.Reply(ctx, evt.RoomID, err.Error()); replyErr != nil {
			return replyErr
		}
		return err
	}

	report, err := h.service.Wakapi.FetchReport(interval)
	if err != nil {
		slog.Error("fetch wakapi report failed", "err", err)
		return err
	}
	if err := h.service.Message.Reply(ctx, evt.RoomID, report); err != nil {
		slog.Error("send wakapi message failed", "report", report, "err", err)
		return err
	}
	return nil
}

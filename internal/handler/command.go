package handler

import (
	"context"
	"log/slog"
	"strings"

	"maunium.net/go/mautrix/event"
)

func (h *Handler) registerCommands() {
	h.commands["!ping"] = h.handlePing
	h.commands["!wakapi"] = h.handleWakapiCommand
	slog.Info("commands registered")
}

func (h *Handler) HandleCommand(ctx context.Context, evt *event.Event) {
	content := evt.Content.AsMessage()
	if content.MsgType != event.MsgText {
		return
	}

	if evt.Timestamp < h.startTime.UnixMilli() {
		slog.Debug("receive old message, skipped", "msg", content.Body)
		return
	}

	if evt.Sender == h.service.Message.UserID() {
		slog.Debug("receive own message, skipped")
		return
	}

	body := strings.TrimSpace(content.Body)
	slog.Debug(
		"received text message",
		"room", evt.RoomID,
		"sender", evt.Sender,
		"body", body,
	)

	for prefix, handler := range h.commands {
		if strings.HasPrefix(body, prefix) {
			if err := handler(ctx, evt); err != nil {
				slog.Error(
					"handle error",
					"err", err,
				)
			}
			return
		}
	}
}

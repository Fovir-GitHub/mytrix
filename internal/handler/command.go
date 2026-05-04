package handler

import (
	"context"
	"log/slog"
	"strings"

	"maunium.net/go/mautrix/event"
)

// registerCommands initializes the command handlers map mapping command prefixes to their respective handler functions.
func (h *Handler) registerCommands() {
	h.commands["!ping"] = h.handlePing
	h.commands["!umami"] = h.handleUmamiCommand
	h.commands["!wakapi"] = h.handleWakapiCommand
	h.commands["!rss"] = h.handleRSSCommand
	h.commands["!version"] = h.handleVersion
	slog.Info("bot command initialized",
		"count", len(h.commands))
}

// HandleCommand processes incoming Matrix message events for bot commands.
// It filters out non-text messages, old messages (before bot start), and messages from the bot itself.
// For valid text messages, it checks if the message starts with any registered command prefix and invokes the corresponding handler.
func (h *Handler) HandleCommand(ctx context.Context, evt *event.Event) {
	content := evt.Content.AsMessage()
	if content.MsgType != event.MsgText {
		return
	}

	if evt.Timestamp < h.startTime.UnixMilli() {
		slog.Debug("message skipped (old timestamp)",
			"room", evt.RoomID.String())
		return
	}

	if evt.Sender == h.service.Message.UserID() {
		slog.Debug("message skipped (own message)")
		return
	}

	body := strings.TrimSpace(content.Body)
	slog.Debug(
		"received text message",
		"room", evt.RoomID,
		"sender", evt.Sender,
		"len", len(body),
	)

	for prefix, handler := range h.commands {
		if strings.HasPrefix(body, prefix) {
			if err := handler(ctx, evt); err != nil {
				slog.Error(
					"command handler failed",
					"room", evt.RoomID.String(),
					"sender", evt.Sender.String(),
					"err", err,
				)
			}
			return
		}
	}
}

package handler

import (
	"context"
	"log/slog"
	"strings"

	"github.com/Fovir-GitHub/mytrix/internal/service"
	"maunium.net/go/mautrix/event"
)

type Handler struct {
	service  *service.Service
	commands map[string]func(context.Context, *event.Event) error
}

func (h *Handler) registerCommands() {
	h.commands["!ping"] = h.handlePing
	h.commands["!gotify"] = h.handleGotify
}

func NewHandler(s *service.Service) *Handler {
	h := &Handler{
		service:  s,
		commands: make(map[string]func(context.Context, *event.Event) error),
	}
	h.registerCommands()

	return h
}

func (h *Handler) Handle(ctx context.Context, evt *event.Event) {
	content := evt.Content.AsMessage()
	if content.MsgType != event.MsgText {
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

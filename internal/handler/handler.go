// Package handler handles incoming events and commands.
package handler

import (
	"context"
	"log/slog"
	"time"

	"codeberg.org/Fovir/mytrix/internal/model"
	"codeberg.org/Fovir/mytrix/internal/service"
	"maunium.net/go/mautrix/event"
)

// Handler manages Matrix event handling and command execution.
// It routes incoming messages to appropriate command handlers and WebSocket event handlers.
type Handler struct {
	service   *service.Service
	commands  map[string]func(context.Context, *event.Event) error
	events    map[string]func(context.Context, *model.WsEvent) error
	startTime time.Time
}

// NewHandler returns a new Handler with the given service.
// It initializes the command and event maps and registers handlers.
func NewHandler(s *service.Service) *Handler {
	h := &Handler{
		service:   s,
		commands:  make(map[string]func(context.Context, *event.Event) error),
		events:    make(map[string]func(context.Context, *model.WsEvent) error),
		startTime: time.Now(),
	}
	h.registerCommands()
	h.registerWSHandler()
	slog.Info("handler initialized")
	return h
}

func (h *Handler) getReply(ctx context.Context, evt *event.Event) func(string) error {
	return func(s string) error {
		return h.service.Message.Reply(ctx, evt.RoomID, s)
	}
}

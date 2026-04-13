package handler

import (
	"context"
	"log/slog"

	"codeberg.org/Fovir/mytrix/internal/model"
)

// registerWSHandler registers WebSocket event handlers mapping event sources to their respective handler functions.
func (h *Handler) registerWSHandler() {
	h.events[model.SourceGotify] = h.handleGotify
	slog.Info("websocket handler registered")
}

// HandleWSEvent handles incoming WebSocket events.
// It routes the event to the appropriate handler based on the event source.
// If no handler is found for the source, it returns nil.
func (h *Handler) HandleWSEvent(ctx context.Context, event *model.WsEvent) error {
	for name, handler := range h.events {
		if event.Source == name {
			slog.Debug("handle event", "name", name)
			return handler(ctx, event)
		}
	}
	return nil
}

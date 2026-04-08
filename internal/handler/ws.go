package handler

import (
	"context"
	"log/slog"

	"github.com/Fovir-GitHub/mytrix/internal/model"
)

func (h *Handler) registerWSHandler() {
	h.events[model.SourceGotify] = h.handleGotify
	slog.Info("websocket handler registered")
}

func (h *Handler) HandleWSEvent(ctx context.Context, event *model.WsEvent) error {
	for name, handler := range h.events {
		if event.Source == name {
			slog.Debug("handle event", "name", name)
			return handler(ctx, event)
		}
	}
	return nil
}

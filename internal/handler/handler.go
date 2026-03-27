package handler

import (
	"context"
	"log/slog"

	"github.com/Fovir-GitHub/mytrix/internal/model"
	"github.com/Fovir-GitHub/mytrix/internal/service"
	"maunium.net/go/mautrix/event"
)

type Handler struct {
	service  *service.Service
	commands map[string]func(context.Context, *event.Event) error
	events   map[string]func(context.Context, *model.WsEvent) error
}

func NewHandler(s *service.Service) *Handler {
	slog.Debug("create handler")
	h := &Handler{
		service:  s,
		commands: make(map[string]func(context.Context, *event.Event) error),
		events:   make(map[string]func(context.Context, *model.WsEvent) error),
	}

	h.registerCommands()
	h.registerWSHandler()
	return h
}

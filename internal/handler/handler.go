package handler

import (
	"context"

	"github.com/Fovir-GitHub/mytrix/internal/service"
	"maunium.net/go/mautrix/event"
)

type Handler struct {
	service  *service.Service
	commands map[string]func(context.Context, *event.Event) error
}

func NewHandler(s *service.Service) *Handler {
	h := &Handler{
		service:  s,
		commands: make(map[string]func(context.Context, *event.Event) error),
	}
	h.registerCommands()

	return h
}

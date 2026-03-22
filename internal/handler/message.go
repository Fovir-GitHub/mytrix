package handler

import (
	"context"

	"github.com/Fovir-GitHub/mytrix/internal/service"
	"maunium.net/go/mautrix/event"
)

type MessageHandler struct {
	service *service.MessageService
}

func NewMessageHandler(s *service.MessageService) *MessageHandler {
	return &MessageHandler{service: s}
}

func (h *MessageHandler) Handle(ctx context.Context, evt *event.Event) {
	h.service.HandleMessage(ctx, evt)
}

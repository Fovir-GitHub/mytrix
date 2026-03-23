package service

import (
	"context"

	"github.com/Fovir-GitHub/mytrix/internal/client"
	"maunium.net/go/mautrix/event"
)

type MessageService struct {
	client *client.MatrixClient
}

func newMessageService(c *client.MatrixClient) *MessageService {
	return &MessageService{client: c}
}

func (s *MessageService) Reply(ctx context.Context, evt *event.Event, text string) error {
	return s.client.SendTextMessage(ctx, evt.RoomID, text)
}

func (s *MessageService) HandleMessage(ctx context.Context, evt *event.Event) {
	content := evt.Content.AsMessage()
	if content.MsgType != event.MsgText || evt.Sender == s.client.UserID() {
		return
	}

	if content.Body == "!ping" {
		_ = s.client.SendTextMessage(ctx, evt.RoomID, "pong")
		return
	}

	_ = s.client.SendTextMessage(ctx, evt.RoomID, "Feedback")
}

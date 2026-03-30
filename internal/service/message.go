package service

import (
	"context"
	"fmt"

	"github.com/Fovir-GitHub/mytrix/internal/matrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type MessageService struct {
	client *matrix.Client
}

func newMessageService(c *matrix.Client) *MessageService {
	return &MessageService{client: c}
}

func (s *MessageService) Reply(ctx context.Context, roomID id.RoomID, text string) error {
	return s.client.SendTextMessage(ctx, roomID, text)
}

func (s *MessageService) Ping(ctx context.Context, evt *event.Event) error {
	err := s.Reply(ctx, evt.RoomID, "pong")
	if err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}
	return nil
}

func (s *MessageService) UserID() id.UserID {
	return s.client.UserID()
}

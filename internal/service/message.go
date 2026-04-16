package service

import (
	"context"
	"log/slog"

	"codeberg.org/Fovir/mytrix/internal/matrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

// MessageService provides message-related operations on the Matrix client.
type MessageService struct {
	client *matrix.Client
}

// NewMessageService creates a new MessageService with the provided Matrix client.
func NewMessageService(c *matrix.Client) *MessageService {
	slog.Info("create message service")
	return &MessageService{client: c}
}

// Reply sends a text message to the specified room.
func (s *MessageService) Reply(ctx context.Context, roomID id.RoomID, text string) error {
	if err := s.client.SendTextMessage(ctx, roomID, text); err != nil {
		slog.Error("send message failed", "roomID", roomID.String(), "text", text, "err", err)
		return err
	}
	slog.Debug("sent message", "roomID", roomID.String(), "len", len(text))
	return nil
}

// Ping sends a "pong" response to the specified room.
func (s *MessageService) Ping(ctx context.Context, evt *event.Event) error {
	return s.Reply(ctx, evt.RoomID, "pong")
}

// UserID returns the user ID of the Matrix client.
func (s *MessageService) UserID() id.UserID {
	return s.client.UserID()
}

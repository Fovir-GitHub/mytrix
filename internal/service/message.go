package service

import (
	"context"
	"log/slog"

	"codeberg.org/Fovir/mytrix/internal/matrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type MessageService struct {
	client *matrix.Client
}

func newMessageService(c *matrix.Client) *MessageService {
	slog.Info("create message service")
	return &MessageService{client: c}
}

func (s *MessageService) Reply(ctx context.Context, roomID id.RoomID, text string) error {
	if err := s.client.SendTextMessage(ctx, roomID, text); err != nil {
		slog.Error("send message failed", "roomID", roomID.String(), "text", text, "err", err)
		return err
	}
	slog.Debug("sent message", "roomID", roomID.String(), "msg", text)
	return nil
}

func (s *MessageService) Ping(ctx context.Context, evt *event.Event) error {
	slog.Debug("ping command called")
	return s.Reply(ctx, evt.RoomID, "pong")
}

func (s *MessageService) UserID() id.UserID {
	return s.client.UserID()
}

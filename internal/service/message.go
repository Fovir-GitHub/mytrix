package service

import (
	"context"
	"log/slog"

	"codeberg.org/Fovir/mytrix/internal/matrix"
	"codeberg.org/Fovir/mytrix/internal/version"
	"maunium.net/go/mautrix"
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

// ReplyWithoutResp sends a text message to the specified room without returning `*mautrix.RespSendEvent`
func (s *MessageService) ReplyWithoutResp(ctx context.Context, roomID id.RoomID, text string) error {
	_, err := s.client.SendTextMessage(ctx, roomID, text)
	if err != nil {
		slog.Error("send message failed", "roomID", roomID.String(), "text", text, "err", err)
		return err
	}
	slog.Debug("sent message", "roomID", roomID.String(), "len", len(text))
	return nil
}

func (s *MessageService) Reply(ctx context.Context, roomID id.RoomID, text string) (*mautrix.RespSendEvent, error) {
	resp, err := s.client.SendTextMessage(ctx, roomID, text)
	if err != nil {
		slog.Error("send message failed", "roomID", roomID.String(), "text", text, "err", err)
		return nil, err
	}
	slog.Debug("sent message", "roomID", roomID.String(), "len", len(text))
	return resp, nil
}

func (s *MessageService) ReplyThread(ctx context.Context, roomID id.RoomID, eventID id.EventID, text string) (*mautrix.RespSendEvent, error) {
	content := &event.MessageEventContent{
		MsgType: event.MsgText,
		Body:    text,
	}

	content.SetThread(&event.Event{
		ID:     eventID,
		RoomID: roomID,
	})

	return s.client.SendEventMessage(ctx, roomID, content)
}

// Ping sends a "pong" response to the specified room.
func (s *MessageService) Ping(ctx context.Context, evt *event.Event) error {
	return s.ReplyWithoutResp(ctx, evt.RoomID, "pong")
}

// UserID returns the user ID of the Matrix client.
func (s *MessageService) UserID() id.UserID {
	return s.client.UserID()
}

// Version sends current version of mytrix.
func (s *MessageService) Version(ctx context.Context, evt *event.Event) error {
	return s.ReplyWithoutResp(ctx, evt.RoomID, version.Version)
}

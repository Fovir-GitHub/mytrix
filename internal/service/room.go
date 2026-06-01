package service

import (
	"context"
	"fmt"
	"log/slog"

	"codeberg.org/Fovir/mytrix/internal/config"
	"codeberg.org/Fovir/mytrix/internal/matrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type RoomService struct {
	client *matrix.Client
}

func NewRoomService(c *matrix.Client) *RoomService {
	slog.Info("create room service")
	return &RoomService{
		client: c,
	}
}

func (r *RoomService) ProcessInvitation(ctx context.Context, evt *event.Event) error {
	admin := config.Config.AdminID
	if len(admin) <= 0 {
		slog.Warn("invitation ignored (admin is not configured)", "sender", evt.Sender)
		return nil
	}

	// Check whether the invitation is from the bot self.
	if evt.Sender == r.client.UserID() {
		slog.Debug("invitation skipped (own invitation)")
		return nil
	}

	// Check invitation target (ignore errors).
	if evt.GetStateKey() != r.client.UserID().String() && evt.Content.AsMember().Membership != event.MembershipInvite {
		return nil
	}

	// Check the inviter is admin or not.
	if evt.Sender != id.UserID(admin) {
		if err := r.client.LeaveRoom(ctx, evt.RoomID); err != nil {
			return fmt.Errorf("inviter is not admin: %v, and leave room failed (id=%v): %w", evt.Sender, evt.RoomID.String(), err)
		}
		return fmt.Errorf("inviter is not admin: %v", evt.Sender)
	}

	// Accept invitation.
	if err := r.client.JoinRoomByID(ctx, evt.RoomID); err != nil {
		return fmt.Errorf("join room failed (id=%v): %w", evt.RoomID.String(), err)
	}

	return nil
}

func (r *RoomService) MaxMessageLength(ctx context.Context, roomID id.RoomID) int {
	return r.client.MaxMessageLength(ctx, roomID)
}

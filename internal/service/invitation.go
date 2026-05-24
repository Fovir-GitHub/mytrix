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

type InvitationService struct {
	client *matrix.Client
}

func NewInvitationService(c *matrix.Client) *InvitationService {
	slog.Info("create invitation service")
	return &InvitationService{client: c}
}

func (i *InvitationService) Process(ctx context.Context, evt *event.Event) error {
	admin := config.Config.AdminID

	// Check invitation target (ignore errors).
	if evt.GetStateKey() != i.client.UserID().String() && evt.Content.AsMember().Membership != event.MembershipInvite {
		return nil
	}

	// Check the inviter is admin or not.
	if evt.Sender != id.UserID(admin) {
		if err := i.client.LeaveRoom(ctx, evt.RoomID); err != nil {
			return fmt.Errorf("inviter is not admin: %v, and leave room failed (id=%v): %w", evt.Sender, evt.RoomID.String(), err)
		}

		return fmt.Errorf("inviter is not admin: %v", evt.Sender)
	}

	// Accept invitation.
	if err := i.client.JoinRoomByID(ctx, evt.RoomID); err != nil {
		return fmt.Errorf("join room failed (id=%v): %w", evt.RoomID.String(), err)
	}

	return nil
}

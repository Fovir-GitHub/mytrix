package service

import (
	"context"
	"fmt"
	"log/slog"

	"codeberg.org/Fovir/mytrix/internal/db"
	"codeberg.org/Fovir/mytrix/internal/matrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type RoomService struct {
	client *matrix.Client
	q      *db.Queries
}

func NewRoomService(c *matrix.Client, q *db.Queries) *RoomService {
	slog.Info("create room service")
	return &RoomService{
		client: c,
		q:      q,
	}
}

func (r *RoomService) ProcessInvitation(ctx context.Context, evt *event.Event) error {
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
	senderIDStr := evt.Sender.String()
	roomIDStr := evt.RoomID.String()
	isAdmin, err := r.q.IsUserAdmin(ctx, senderIDStr)
	if err != nil {
		return fmt.Errorf("check user %v is admin failed: %w", senderIDStr, err)
	}

	// If the inviter is not admin, then the bot should leave the room, set the room state to `left`.
	if !isAdmin {
		if err := r.client.LeaveRoom(ctx, evt.RoomID); err != nil {
			return fmt.Errorf("inviter is not admin: %v, and leave room failed (id=%v): %w", evt.Sender, evt.RoomID.String(), err)
		}

		if err := r.updateRoomState(ctx, roomIDStr, "left"); err != nil {
			return fmt.Errorf("set leave room db state failed (room_id=%v): %w", roomIDStr, err)
		}

		return fmt.Errorf("inviter is not admin: %v", evt.Sender)
	}

	// Accept invitation, and set the room state to `joined`.
	if err := r.client.JoinRoomByID(ctx, evt.RoomID); err != nil {
		return fmt.Errorf("join room failed (id=%v): %w", evt.RoomID.String(), err)
	}

	if err := r.updateRoomState(ctx, roomIDStr, "joined"); err != nil {
		return fmt.Errorf("set join room db state fialed (room_id=%v): %w", roomIDStr, err)
	}

	return nil
}

func (r *RoomService) MaxMessageLength(ctx context.Context, roomID id.RoomID) int {
	return r.client.MaxMessageLength(ctx, roomID)
}

func (r *RoomService) updateRoomState(ctx context.Context, roomID, state string) error {
	if err := r.q.CreateRoom(ctx, &db.CreateRoomParams{
		ID:    roomID,
		State: state,
	}); err != nil {
		return err
	}

	return r.q.UpdateRoomState(ctx, &db.UpdateRoomStateParams{
		ID:    roomID,
		State: state,
	})
}

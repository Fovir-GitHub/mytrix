package handler

import (
	"context"
	"log/slog"
	"strings"

	"maunium.net/go/mautrix/event"
)

func (h *Handler) handleUserCommand(ctx context.Context, evt *event.Event) error {
	msg := evt.Content.AsMessage().Body
	parts := strings.Fields(msg)
	if len(parts) <= 1 {
		return h.handleUserHelp(ctx, evt)
	}

	switch parts[1] {
	case "list":
		return h.handleUserList(ctx, evt)
	default:
		return h.handleUserHelp(ctx, evt)
	}
}

func (h *Handler) handleUserList(ctx context.Context, evt *event.Event) error {
	reply := h.getReply(ctx, evt)
	users, err := h.service.User.ListUsers(ctx)
	if err != nil {
		slog.Error("list users failed", "err", err)
		return reply("Failed to list users")
	}

	if users == "" {
		return reply("No users")
	}

	return reply(users)
}

func (h *Handler) handleUserHelp(ctx context.Context, evt *event.Event) error {
	const userCommandUsage = "Usage:\n" +
		"```" + `
!user list		List all users
` + "```"
	reply := h.getReply(ctx, evt)
	return reply(userCommandUsage)
}

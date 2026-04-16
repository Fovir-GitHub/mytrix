// TODO: Compose the help information
package handler

import (
	"context"
	"log/slog"
	"strconv"
	"strings"

	"maunium.net/go/mautrix/event"
)

// handleRSSCommand processes the !rss command with various subcommands (add, delete, list).
func (h *Handler) handleRSSCommand(ctx context.Context, evt *event.Event) error {
	msg := evt.Content.AsMessage().Body
	parts := strings.Fields(msg)
	if len(parts) <= 1 {
		return h.service.Message.Reply(ctx, evt.RoomID, "invalid argument")
	}

	switch parts[1] {
	case "add":
		return h.handleRSSAdd(ctx, evt, parts)
	case "delete":
		return h.handleRSSDelete(ctx, evt, parts)
	case "list":
		return h.handleRSSList(ctx, evt)
	default:
		return h.service.Message.Reply(ctx, evt.RoomID, "invalid argument")
	}
}

// handleRSSAdd adds a new RSS feed to the subscription list.
func (h *Handler) handleRSSAdd(ctx context.Context, evt *event.Event, parts []string) error {
	reply := h.getReply(ctx, evt)
	if len(parts) < 3 {
		return reply("invalid arguments")
	}

	u := parts[2]
	if err := h.service.RSS.AddFeed(u); err != nil {
		slog.Error("add rss failed", "url", u, "err", err)
		return reply("failed to add RSS feed")
	}

	h.handleRSSSchedule(ctx)
	return reply("RSS feed added successfully")
}

// handleRSSDelete deletes a RSS feed.
func (h *Handler) handleRSSDelete(ctx context.Context, evt *event.Event, parts []string) error {
	reply := h.getReply(ctx, evt)

	if len(parts) < 3 {
		return reply("invalid arguments")
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		return reply("invalid id")
	}
	if err := h.service.RSS.DeleteFeed(id); err != nil {
		slog.Error("delete rss feed failed", "id", id, "err", err)
		return reply("failed to delete RSS feed")
	}
	return reply("feed deleted")
}

// handleRSSList lists all RSS feeds.
func (h *Handler) handleRSSList(ctx context.Context, evt *event.Event) error {
	reply := h.getReply(ctx, evt)
	feeds, err := h.service.RSS.ListFeeds()
	if err != nil {
		slog.Error("list RSS feeds failed", "err", err)
		return reply("failed to list RSS feeds")
	}
	if feeds == "" {
		return reply("Empty RSS list")
	}

	return reply(feeds)
}

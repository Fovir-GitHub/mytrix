// TODO: Compose the help information
package handler

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"maunium.net/go/mautrix/event"
)

func (h *Handler) handleRSSCommand(ctx context.Context, evt *event.Event) error {
	msg := evt.Content.AsMessage().Body
	parts := strings.Fields(msg)
	if len(parts) <= 2 {
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

func (h *Handler) handleRSSAdd(ctx context.Context, evt *event.Event, parts []string) error {
	reply := h.getReply(ctx, evt)
	if len(parts) < 3 {
		return reply("invalid arguments")
	}

	u := parts[2]
	if err := h.service.RSS.AddFeed(u); err != nil {
		slog.Error("add rss failed", "url", u, "err", err)
		return reply("failed to add RSS feed")
	} else {
		return reply("RSS feed added successfully")
	}
}

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

func (h *Handler) handleRSSList(ctx context.Context, evt *event.Event) error {
	reply := h.getReply(ctx, evt)
	feeds, err := h.service.RSS.AllFeeds()
	if err != nil {
		slog.Error("list RSS feeds failed", "err", err)
		return reply("failed to list RSS feeds")
	}

	if len(feeds) <= 0 {
		return reply("empty list")
	}

	var msg strings.Builder

	// TODO: format the feed information
	for _, feed := range feeds {
		fmt.Fprintf(&msg, "%d %s %s\n", feed.ID, feed.Title, feed.URL)
	}
	return reply(msg.String())
}

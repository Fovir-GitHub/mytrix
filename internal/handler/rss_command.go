package handler

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"strings"

	"codeberg.org/Fovir/mytrix/internal/config"
	"codeberg.org/Fovir/mytrix/internal/service"
	"maunium.net/go/mautrix/event"
)

// handleRSSCommand processes the !rss command with various subcommands (add, delete, list, export, update, help).
func (h *Handler) handleRSSCommand(ctx context.Context, evt *event.Event) error {
	msg := evt.Content.AsMessage().Body
	parts := strings.Fields(msg)
	if len(parts) <= 1 {
		return h.handleRSSHelp(ctx, evt)
	}

	switch parts[1] {
	case "add":
		return h.handleRSSAdd(ctx, evt, parts)
	case "delete":
		return h.handleRSSDelete(ctx, evt, parts)
	case "list":
		return h.handleRSSList(ctx, evt)
	case "export":
		return h.handleRSSExport(ctx, evt)
	case "update":
		return h.handleRSSUpdate(ctx, evt)
	case "help":
		return h.handleRSSHelp(ctx, evt)
	default:
		return h.handleRSSHelp(ctx, evt)
	}
}

// handleRSSAdd adds new RSS feeds to the subscription list.
func (h *Handler) handleRSSAdd(ctx context.Context, evt *event.Event, parts []string) error {
	reply := h.getReply(ctx, evt)
	if len(parts) < 3 {
		return h.handleRSSHelp(ctx, evt)
	}

	feeds := parts[2:]
	var replyMsg string
	errFeeds, err := h.service.RSS.AddFeeds(feeds)

	if err != nil {
		slog.Error("add rss failed", "err", err)
		replyMsg = "Failed to add RSS feeds:\n" + errFeeds
	} else {
		replyMsg = "RSS feeds added successfully"
	}

	if config.Config.RSS.UpdateAfterAdd {
		h.handleRSSSchedule(ctx)
	}

	return reply(replyMsg)
}

// handleRSSDelete deletes a RSS feed.
func (h *Handler) handleRSSDelete(ctx context.Context, evt *event.Event, parts []string) error {
	reply := h.getReply(ctx, evt)

	if len(parts) < 3 {
		return h.handleRSSHelp(ctx, evt)
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		return reply("Invalid ID")
	}
	if err := h.service.RSS.DeleteFeed(id); err != nil {
		slog.Error("delete rss feed failed", "id", id, "err", err)
		return reply("Failed to delete RSS feed")
	}
	return reply("Feed deleted")
}

// handleRSSList lists all RSS feeds.
func (h *Handler) handleRSSList(ctx context.Context, evt *event.Event) error {
	reply := h.getReply(ctx, evt)
	feeds, err := h.service.RSS.ListFeeds()
	if err != nil {
		slog.Error("list RSS feeds failed", "err", err)
		return reply("Failed to list RSS feeds")
	}
	if feeds == "" {
		return reply("Empty RSS list")
	}

	return reply(feeds)
}

// handleRSSExport exports all RSS feeds.
func (h *Handler) handleRSSExport(ctx context.Context, evt *event.Event) error {
	reply := h.getReply(ctx, evt)
	feeds, err := h.service.RSS.ExportFeeds()
	if err != nil {
		slog.Error("export RSS feeds failed", "err", err)
		return reply("Failed to export RSS feeds")
	}
	if feeds == "" {
		return reply("Empty RSS list")
	}
	return reply(feeds)
}

// handleRSSUpdate updates RSS feeds manually.
func (h *Handler) handleRSSUpdate(ctx context.Context, evt *event.Event) error {
	reply := h.getReply(ctx, evt)
	updated, err := h.service.RSS.Update()
	if err != nil {
		if errors.Is(err, service.ErrRSSFetchFeeds) {
			slog.Error("update rss error", "err", err)
			return reply("Failed to update RSS feeds")
		}
		if errors.Is(err, service.ErrRSSNoUpdate) {
			return reply("Everything up to date")
		}
	}
	return reply(updated)
}

// handleRSSHelp shows command help information of RSS.
func (h *Handler) handleRSSHelp(ctx context.Context, evt *event.Event) error {
	const rssCommandUsage = "Usage:\n" +
		"```" + `
!rss add <url1> <url2> ...  Add one or more feeds
!rss delete <id>            Delete a feed
!rss list                   List feeds
!rss export                 Export feeds
!rss update                 Update all feeds
` + "```"
	reply := h.getReply(ctx, evt)
	return reply(rssCommandUsage)
}

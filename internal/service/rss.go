package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"codeberg.org/Fovir/mytrix/internal/config"
	"codeberg.org/Fovir/mytrix/internal/db"
	"codeberg.org/Fovir/mytrix/internal/feed"
	"codeberg.org/Fovir/mytrix/internal/render"
)

type RSSService interface {
	AddFeeds(context.Context, []string) (string, error)
	DeleteFeeds(context.Context, []string) (string, error)
	Update(context.Context) ([]string, error)
	ListFeeds(context.Context) (string, error)
	ExportFeeds(context.Context) (string, error)
}

type RealRSSService struct {
	q      *db.Queries
	parser *feed.Parser
}

func NewRSSService(query *db.Queries) RSSService {
	cfg := config.Config.RSS
	slog.Info("rss service initialized", "enabled", cfg.Enabled)
	if !cfg.Enabled {
		return &NoopRSSService{err: fmt.Errorf("RSS is not enabled")}
	}
	return &RealRSSService{
		q:      query,
		parser: feed.New(),
	}
}

func (r *RealRSSService) AddFeeds(ctx context.Context, feeds []string) (string, error) {
	var errFeeds strings.Builder
	var errs []error

	for _, f := range feeds {
		if err := r.addFeed(ctx, f); err != nil {
			errFeeds.WriteString(f)
			errFeeds.WriteString("\n")
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errFeeds.String(), errors.Join(errs...)
	}
	return "", nil
}

func (r *RealRSSService) addFeed(ctx context.Context, u string) error {
	feed, _, err := r.parser.ParseURL(u)
	if err != nil {
		return fmt.Errorf("parse rss url failed (url=%s): %w", u, err)
	}

	err = r.q.CreateRSSFeed(ctx, &db.CreateRSSFeedParams{Url: feed.Url, Title: feed.Title})
	if err != nil {
		return fmt.Errorf("create rss feed failed (url=%s): %w", u, err)
	}
	slog.Info("rss feed added", "url", u)
	return nil
}

func (r *RealRSSService) DeleteFeeds(ctx context.Context, ids []string) (string, error) {
	var errs []error
	var errIds strings.Builder
	handleErr := func(idStr string, err error) {
		errIds.WriteString(idStr)
		errIds.WriteString("\n")
		errs = append(errs, err)
	}

	for _, idStr := range ids {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			handleErr(idStr, err)
			continue
		}

		if err := r.deleteFeed(ctx, id); err != nil {
			handleErr(idStr, err)
		}
	}

	if len(errs) > 0 {
		return errIds.String(), errors.Join(errs...)
	}
	return "", nil
}

func (r *RealRSSService) deleteFeed(ctx context.Context, id int) error {
	if err := r.q.DeleteRSSFeed(ctx, int64(id)); err != nil {
		return fmt.Errorf("delete feed failed (id=%d): %w", id, err)
	}
	if err := r.q.DeleteRSSItemByFeedID(ctx, int64(id)); err != nil {
		return fmt.Errorf("delete feed items failed (feed_id=%d): %w", id, err)
	}
	slog.Info("rss feed deleted", "id", id)
	return nil
}

func (r *RealRSSService) Update(ctx context.Context) ([]string, error) {
	var (
		errs []error
		res  []string
	)

	feeds, err := r.allFeeds(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRSSFetchFeeds, err)
	}
	slog.Debug("rss update start", "feeds_len", len(feeds))

	for _, feed := range feeds {
		updated, err := r.updateFeed(ctx, &feed)
		if err != nil {
			errs = append(errs, err)
			slog.Warn("feed update failed", "feed_id", feed.ID, "err", err)
		}

		if len(updated) <= 0 {
			continue
		}

		res = append(res, updated)
	}

	if len(res) <= 0 {
		return nil, ErrRSSNoUpdate
	}

	if len(errs) > 0 {
		return res, ErrRSSPartialUpdate
	}

	slog.Info("rss update finished", "feeds_len", len(feeds))

	return res, nil
}

func (r *RealRSSService) updateFeed(ctx context.Context, feed *db.RssFeed) (string, error) {
	var (
		updated strings.Builder
		errs    []error
	)

	_, items, err := r.parser.ParseURL(feed.Url)
	if err != nil {
		return "", err
	}

	for _, item := range items {
		item.FeedID = feed.ID
		if err := r.addItem(ctx, &item); err != nil {
			if !errors.Is(err, ErrRSSItemExists) {
				slog.Error("item insert failed", "feed_url", feed.Url, "guid", item.Guid, "err", err)
				errs = append(errs, fmt.Errorf("insert item failed (feed_url=%s, guid=%s): %w", feed.Url, item.Guid, err))
			} else {
				slog.Debug("item insert failed", "feed_url", feed.Url, "guid", item.Guid, "err", err)
			}
			continue
		}
		updated.WriteString(render.RssItemMarkdown(feed, &item))
		updated.WriteString("\n")
	}
	if len(errs) > 0 {
		slog.Warn(
			"some items failed",
			"feed_url", feed.Url,
			"failed", len(errs),
			"total", len(items),
		)
		return updated.String(), fmt.Errorf("update feed failed (url=%s): %w", feed.Url, errors.Join(errs...))
	}

	return updated.String(), nil
}

func (r *RealRSSService) addItem(ctx context.Context, item *db.RssItem) error {
	err := r.q.CreateRSSItem(ctx, &db.CreateRSSItemParams{
		FeedID:      item.FeedID,
		Guid:        item.Guid,
		Link:        item.Link,
		Title:       item.Title,
		Description: item.Description,
	})
	if err != nil {
		return fmt.Errorf("add item failed (feed_id=%d, guid=%s): %w", item.FeedID, item.Guid, err)
	}
	return nil
}

func (r *RealRSSService) allFeeds(ctx context.Context) ([]db.RssFeed, error) {
	feeds, err := r.q.AllFeeds(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch all feeds failed: %w", err)
	}
	return feeds, nil
}

func (r *RealRSSService) ListFeeds(ctx context.Context) (string, error) {
	var res strings.Builder
	feeds, err := r.allFeeds(ctx)
	if err != nil {
		return "", fmt.Errorf("list feed failed: %w", err)
	}
	slog.Debug("list rss feeds", "feeds", len(feeds))
	if len(feeds) <= 0 {
		return "", nil
	}

	for _, feed := range feeds {
		res.WriteString(render.RssFeedMarkdown(&feed))
		res.WriteString("\n")
	}
	return res.String(), nil
}

func (r *RealRSSService) ExportFeeds(ctx context.Context) (string, error) {
	feeds, err := r.allFeeds(ctx)
	if err != nil {
		return "", fmt.Errorf("export feed failed: %w", err)
	}
	var res strings.Builder
	for _, feed := range feeds {
		res.WriteString(feed.Url)
		res.WriteString("\n")
	}
	return res.String(), nil
}

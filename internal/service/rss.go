package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"codeberg.org/Fovir/mytrix/internal/config"
	"codeberg.org/Fovir/mytrix/internal/db"
	"codeberg.org/Fovir/mytrix/internal/feed"
	"codeberg.org/Fovir/mytrix/internal/model"
	"codeberg.org/Fovir/mytrix/internal/render"
	"codeberg.org/Fovir/mytrix/internal/utils"
	"maunium.net/go/mautrix/id"
)

type RSSService interface {
	AddFeeds(context.Context, []string) (string, error)
	DeleteFeeds(context.Context, []string) (string, error)
	Update(context.Context) ([]model.RSSUpdateResult, error)
	ListFeeds(context.Context) (string, error)
	ExportFeeds(context.Context) (string, error)
	MarkItemsPushed(context.Context, *model.RSSUpdateResult) error
	SetThreadRoot(ctx context.Context, feedID int64, roomID id.RoomID, eventID id.EventID) error
}

type RealRSSService struct {
	db     *sql.DB
	q      *db.Queries
	parser *feed.Parser
}

func NewRSSService(sqlDB *sql.DB, query *db.Queries) RSSService {
	cfg := config.Config.RSS
	slog.Info("rss service initialized", "enabled", cfg.Enabled)
	if !cfg.Enabled {
		return &NoopRSSService{err: fmt.Errorf("RSS is not enabled")}
	}
	return &RealRSSService{
		db:     sqlDB,
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

func (r *RealRSSService) Update(ctx context.Context) ([]model.RSSUpdateResult, error) {
	var (
		errs []error
		res  []model.RSSUpdateResult
	)

	// Query all feeds.
	feeds, err := r.allFeeds(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRSSFetchFeeds, err)
	}
	slog.Debug("rss update start", "feeds_len", len(feeds))

	// Update feed by feed.
	for _, feed := range feeds {
		rssUpdateResult, err := r.updateFeed(ctx, &feed)
		if err != nil {
			errs = append(errs, err)
			slog.Warn("feed update failed", "feed_id", feed.ID, "err", err)
			continue
		}

		// Feed is up to date.
		if rssUpdateResult == nil {
			continue
		}

		// Prepend feed title to rendered markdown content.
		rssUpdateResult.Rendered = fmt.Sprintf("# %s\n%s", feed.Title, rssUpdateResult.Rendered)

		res = append(res, *rssUpdateResult)
	}

	if len(errs) > 0 {
		return res, fmt.Errorf("%w: %w", ErrRSSPartialUpdate, errors.Join(errs...))
	}

	if len(res) <= 0 {
		return nil, ErrRSSNoUpdate
	}

	slog.Info("rss update finished", "feeds_len", len(feeds))

	return res, nil
}

func (r *RealRSSService) updateFeed(ctx context.Context, feed *db.RSSFeed) (_ *model.RSSUpdateResult, err error) {
	var rendered strings.Builder
	var ids []int64

	_, items, err := r.parser.ParseURL(feed.Url)
	if err != nil {
		return nil, err
	}

	// Create a transaction.
	tx, qtx, err := utils.CreateTransaction(ctx, r.db, r.q, nil)
	if err != nil {
		return nil, fmt.Errorf("create tx for updating rss feed failed: %w", err)
	}
	defer func() {
		if err != nil {
			if txErr := tx.Rollback(); txErr != nil {
				slog.Error("tx rollback failed", "err", err)
			}
		}
	}()

	// Handle unpushed items.
	unpushed, qErr := qtx.UnpushedRSSItemsByFeedID(ctx, feed.ID)
	if qErr != nil {
		slog.Warn("query unpushed rss item failed", "feed_id", feed.ID, "err", err)
	}

	for _, item := range unpushed {
		rendered.WriteString(render.RSSItemMarkdown(feed, &item))
		rendered.WriteString("\n")
		ids = append(ids, item.ID)
	}

	for _, item := range items {
		item.FeedID = feed.ID
		itemID, insErr := qtx.CreateRSSItem(ctx, &db.CreateRSSItemParams{
			FeedID:      item.FeedID,
			Guid:        item.Guid,
			Link:        item.Link,
			Title:       item.Title,
			Description: item.Description,
		})
		if insErr != nil {
			// Item exists.
			if errors.Is(insErr, sql.ErrNoRows) {
				continue
			}

			err = fmt.Errorf("insert rss item failed (feed_url=%v, guid=%v): %w", feed.Url, item.Guid, err)
			return nil, err
		}

		ids = append(ids, itemID)
		rendered.WriteString(render.RSSItemMarkdown(feed, &item))
		rendered.WriteString("\n")
	}

	if len(ids) <= 0 {
		_ = tx.Rollback()
		return nil, nil
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	if feed.EventID == nil {
		return &model.RSSUpdateResult{
			Rendered: rendered.String(),
			ItemIDs:  ids,
			Event:    nil,
			FeedID:   feed.ID,
		}, nil
	}

	event, err := r.q.GetFeedThreadRootEventByID(ctx, feed.ID)
	if err != nil {
		return nil, err
	}

	return &model.RSSUpdateResult{
		Rendered: rendered.String(),
		ItemIDs:  ids,
		Event:    &event,
		FeedID:   feed.ID,
	}, nil
}

func (r *RealRSSService) allFeeds(ctx context.Context) ([]db.RSSFeed, error) {
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
		res.WriteString(render.RSSFeedMarkdown(&feed))
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

func (r *RealRSSService) MarkItemsPushed(ctx context.Context, rssUpdateResults *model.RSSUpdateResult) (err error) {
	ids := rssUpdateResults.ItemIDs
	tx, qtx, err := utils.CreateTransaction(ctx, r.db, r.q, nil)
	if err != nil {
		return fmt.Errorf("create tx for marking items pushed failed: %w", err)
	}
	defer func() {
		if err != nil {
			if txErr := tx.Rollback(); txErr != nil {
				slog.Error("tx rollback failed", "err", txErr)
			}
		}
	}()

	for _, id := range ids {
		if err = qtx.MarkRSSItemPushedByID(ctx, id); err != nil {
			return fmt.Errorf("mark rss item pushed failed (item_id=%v): %w", id, err)
		}
	}

	return tx.Commit()
}

func (r *RealRSSService) SetThreadRoot(ctx context.Context, feedID int64, roomID id.RoomID, eventID id.EventID) error {
	eventExists := true
	threadRootEvent, err := r.q.GetEventByRoomIDAndEventID(ctx, &db.GetEventByRoomIDAndEventIDParams{
		RoomID:  roomID.String(),
		EventID: eventID.String(),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			eventExists = false
		} else {
			return err
		}
	}

	if eventExists {
		return r.q.UpdateEventByID(ctx, &db.UpdateEventByIDParams{
			RoomID:  roomID.String(),
			EventID: eventID.String(),
			ID:      threadRootEvent.ID,
		})
	}

	newEventID, err := r.q.CreateEvent(ctx, &db.CreateEventParams{
		EventID: eventID.String(),
		RoomID:  roomID.String(),
	})
	if err != nil {
		return err
	}

	return r.q.UpdateFeedThreadRootEventByID(ctx, &db.UpdateFeedThreadRootEventByIDParams{
		EventID: &newEventID,
		ID:      feedID,
	})
}

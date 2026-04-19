package service

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"codeberg.org/Fovir/mytrix/internal/config"
	"codeberg.org/Fovir/mytrix/internal/feed"
	"codeberg.org/Fovir/mytrix/internal/model"
	"codeberg.org/Fovir/mytrix/internal/repo"
	"gorm.io/gorm"
)

type RSSService interface {
	AddFeed(u string) error
	DeleteFeed(id int) error
	Update() (string, error)
	ListFeeds() (string, error)
	ExportFeeds() (string, error)
}

type RealRSSService struct {
	feedRepo *repo.RSSFeedRepo
	itemRepo *repo.RSSItemRepo
	parser   *feed.Parser
}

func NewRSSService(db *gorm.DB) RSSService {
	cfg := config.Config.RSS
	slog.Info("rss service initialized", "enabled", cfg.Enabled)
	if !cfg.Enabled {
		return &NoopRSSService{
			err: fmt.Errorf("RSS is not enabled"),
		}
	}
	feedRepo := repo.NewRSSFeedRepo(db)
	itemRepo := repo.NewRSSItemRepo(db)
	return &RealRSSService{
		feedRepo: feedRepo,
		itemRepo: itemRepo,
		parser:   feed.New(),
	}
}

func (r *RealRSSService) AddFeed(u string) error {
	feed, _, err := r.parser.ParseURL(u)
	if err != nil {
		return fmt.Errorf("parse rss url failed (url=%s): %w", u, err)
	}
	if err := r.feedRepo.Create(feed); err != nil {
		return fmt.Errorf("create rss feed failed (url=%s): %w", u, err)
	}
	slog.Info("rss feed added", "url", u)
	return nil
}

func (r *RealRSSService) DeleteFeed(id int) error {
	if err := r.feedRepo.Delete(id); err != nil {
		return fmt.Errorf("delete feed failed (id=%d): %w", id, err)
	}
	if err := r.itemRepo.DeleteByFeedId(id); err != nil {
		return fmt.Errorf("delete feed items failed (feed_id=%d): %w", id, err)
	}
	slog.Info("rss feed deleted", "id", id)
	return nil
}

func (r *RealRSSService) Update() (string, error) {
	var (
		errs []error
		res  strings.Builder
	)

	feeds, err := r.allFeeds()
	if err != nil {
		return "", fmt.Errorf("update feeds failed: %w", err)
	}
	slog.Debug("rss update start", "feeds_len", len(feeds))

	for _, feed := range feeds {
		updated, err := r.updateFeed(&feed)
		if err != nil {
			errs = append(errs, err)
			slog.Warn("feed update failed", "feed_id", feed.ID, "err", err)
		}
		res.WriteString(updated)
	}
	if len(errs) > 0 {
		return res.String(), errors.Join(errs...)
	}
	slog.Info("rss update finished", "feeds_len", len(feeds))

	return res.String(), nil
}

func (r *RealRSSService) updateFeed(feed *model.RSSFeed) (string, error) {
	var (
		updated strings.Builder
		errs    []error
	)

	_, items, err := r.parser.ParseURL(feed.URL)
	if err != nil {
		return "", err
	}

	for _, item := range items {
		item.FeedID = feed.ID
		if err := r.addItem(&item); err != nil {
			slog.Debug("item insert failed", "feed_url", feed.URL, "guid", item.GUID)
			errs = append(errs, fmt.Errorf("insert item failed (feed_url=%s, guid=%s): %w", feed.URL, item.GUID, err))
			continue
		}
		updated.WriteString(item.ToMarkdown())
		updated.WriteString("\n")
	}
	if len(errs) > 0 {
		slog.Warn(
			"some items failed",
			"feed_url", feed.URL,
			"failed", len(errs),
			"total", len(items),
		)
		return updated.String(), fmt.Errorf("update feed failed (url=%s): %w", feed.URL, errors.Join(errs...))
	}

	return updated.String(), nil
}

func (r *RealRSSService) addItem(item *model.RSSItem) error {
	if err := r.itemRepo.Create(item); err != nil {
		return fmt.Errorf("add item failed (feed_id=%d, guid=%s): %w", item.FeedID, item.GUID, err)
	}
	return nil
}

func (r *RealRSSService) allFeeds() ([]model.RSSFeed, error) {
	feeds, err := r.feedRepo.AllFeeds()
	if err != nil {
		return nil, fmt.Errorf("fetch all feeds failed: %w", err)
	}
	return feeds, nil
}

func (r *RealRSSService) ListFeeds() (string, error) {
	var res strings.Builder
	feeds, err := r.allFeeds()
	if err != nil {
		return "", fmt.Errorf("list feed failed: %w", err)
	}
	slog.Debug("list rss feeds", "feeds", len(feeds))
	if len(feeds) <= 0 {
		return "", nil
	}

	for _, feed := range feeds {
		res.WriteString(feed.ToMarkdown())
		res.WriteString("\n")
	}
	return res.String(), nil
}

func (r *RealRSSService) ExportFeeds() (string, error) {
	feeds, err := r.allFeeds()
	if err != nil {
		return "", fmt.Errorf("export feed failed: %w", err)
	}
	var res strings.Builder
	for _, feed := range feeds {
		res.WriteString(feed.URL)
		res.WriteString("\n")
	}
	return res.String(), nil
}

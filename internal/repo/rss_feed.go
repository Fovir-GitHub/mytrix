package repo

import (
	"errors"
	"fmt"
	"log/slog"

	"codeberg.org/Fovir/mytrix/internal/model"
	"gorm.io/gorm"
)

// RSSFeedRepo provides database operations for RSSFeed entities.
// It wraps a GORM database connection for feed persistence.
type RSSFeedRepo struct {
	db *gorm.DB
}

// NewRSSFeedRepo creates a new RSSFeedRepository with the given database connection.
// It returns a pointer to the initialized repository.
func NewRSSFeedRepo(db *gorm.DB) *RSSFeedRepo {
	return &RSSFeedRepo{db: db}
}

// Create persists the given RSSFeed to the database.
// It uses GORM's Create method and returns any error encountered.
func (r *RSSFeedRepo) Create(feed *model.RSSFeed) error {
	var existing model.RSSFeed
	err := r.db.Unscoped().Where("url = ?", feed.URL).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := r.db.Create(feed).Error; err != nil {
			return fmt.Errorf("create rss feed failed (url=%s): %w", feed.URL, err)
		}
		return nil
	}
	if err != nil {
		return fmt.Errorf("query existing feed failed (url=%s): %w", feed.URL, err)
	}

	if existing.DeletedAt.Valid {
		existing.DeletedAt = gorm.DeletedAt{}
		if err := r.db.Unscoped().Save(&existing).Error; err != nil {
			return fmt.Errorf("save rss feed failed (url=%s): %w", feed.URL, err)
		}
	}

	return nil
}

// Delete removes an RSSFeed from the database by ID.
func (r *RSSFeedRepo) Delete(id int) error {
	if err := r.db.Delete(&model.RSSFeed{}, id).Error; err != nil {
		return fmt.Errorf("delete rss feed failed (id=%d): %w", id, err)
	}
	return nil
}

// SelectFeedByID retrieves a specific RSSFeed from the database by its ID.
func (r *RSSFeedRepo) SelectFeedByID(id int) (*model.RSSFeed, error) {
	var feed *model.RSSFeed
	if err := r.db.First(&feed, id).Error; err != nil {
		return nil, fmt.Errorf("select feed by id failed (id=%d): %w", id, err)
	}
	return feed, nil
}

// AllFeeds retrieves all RSSFeeds from the database.
func (r *RSSFeedRepo) AllFeeds() ([]model.RSSFeed, error) {
	var feeds []model.RSSFeed
	if err := r.db.Find(&feeds).Error; err != nil {
		return nil, fmt.Errorf("get all feeds failed: %w", err)
	}
	slog.Debug("fetch feeds", "len", len(feeds))
	return feeds, nil
}

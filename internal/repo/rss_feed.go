package repo

import (
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
	return r.db.Create(feed).Error
}

func (r *RSSFeedRepo) Delete(id int) error {
	return r.db.Delete(&model.RSSFeed{}, id).Error
}

func (r *RSSFeedRepo) SelectFeedByID(id int) (*model.RSSFeed, error) {
	var feed *model.RSSFeed
	if err := r.db.First(&feed, id).Error; err != nil {
		return nil, err
	}
	return feed, nil
}

func (r *RSSFeedRepo) AllFeeds() ([]model.RSSFeed, error) {
	var feeds []model.RSSFeed
	if err := r.db.Find(&feeds).Error; err != nil {
		return nil, err
	}
	slog.Debug("fetch feeds", "len", len(feeds))
	return feeds, nil
}

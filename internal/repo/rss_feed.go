package repo

import (
	"github.com/Fovir-GitHub/mytrix/internal/model"
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

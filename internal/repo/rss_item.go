package repo

import (
	"errors"
	"fmt"

	"codeberg.org/Fovir/mytrix/internal/model"
	"gorm.io/gorm"
)

// RSSItemRepo provides database operations for RSSItem entities.
// It wraps a GORM database connection for item persistence.
type RSSItemRepo struct {
	db *gorm.DB
}

// NewRSSItemRepo creates a new RSSItemRepository with the given database connection.
// It returns a pointer to the initialized repository.
func NewRSSItemRepo(db *gorm.DB) *RSSItemRepo {
	return &RSSItemRepo{db: db}
}

// Create persists the given RSSItem to the database.
// It uses GORM's Create method and returns any error encountered.
// It also checks items soft-deleted before.
// If the item to be added exists but soft-deleted, it will restore it.
func (r *RSSItemRepo) Create(item *model.RSSItem) error {
	var existing model.RSSItem
	err := r.db.Unscoped().Where("feed_id = ? AND guid = ?", item.FeedID, item.GUID).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := r.db.Create(item).Error; err != nil {
			return fmt.Errorf("create feed item failed (guid=%s): %w", item.GUID, err)
		}
		return nil
	}
	if err != nil {
		return fmt.Errorf("query existing rss item failed (guid=%s): %w", item.GUID, err)
	}

	if existing.DeletedAt.Valid {
		existing.DeletedAt = gorm.DeletedAt{}
		if err := r.db.Unscoped().Save(&existing).Error; err != nil {
			return fmt.Errorf("save rss item failed (guid=%s): %w", item.GUID, err)
		}
	}
	return nil
}

func (r *RSSItemRepo) DeleteByFeedId(feedId int) error {
	if err := r.db.Where("feed_id = ?", feedId).Delete(&model.RSSItem{}).Error; err != nil {
		return fmt.Errorf("delete rss item failed (feed_id=%d), %w", feedId, err)
	}
	return nil
}

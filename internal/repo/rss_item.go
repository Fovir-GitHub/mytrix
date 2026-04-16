package repo

import (
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
func (r *RSSItemRepo) Create(item *model.RSSItem) error {
	if err := r.db.Create(item).Error; err != nil {
		return fmt.Errorf("create feed item failed (id=%d): %w", item.ID, err)
	}
	return nil
}

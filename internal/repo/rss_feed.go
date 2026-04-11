package repo

import (
	"github.com/Fovir-GitHub/mytrix/internal/model"
	"gorm.io/gorm"
)

type RSSFeedRepo struct {
	db *gorm.DB
}

func NewRSSFeedRepo(db *gorm.DB) *RSSFeedRepo {
	return &RSSFeedRepo{db: db}
}

func (r *RSSFeedRepo) Create(feed *model.RSSFeed) error {
	return r.db.Create(feed).Error
}

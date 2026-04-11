package repo

import (
	"github.com/Fovir-GitHub/mytrix/internal/model"
	"gorm.io/gorm"
)

type RSSItemRepo struct {
	db *gorm.DB
}

func NewRSSItemRepo(db *gorm.DB) *RSSItemRepo {
	return &RSSItemRepo{db: db}
}

func (r *RSSItemRepo) Create(item *model.RSSItem) error {
	return r.db.Create(item).Error
}

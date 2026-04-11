package database

import (
	"github.com/Fovir-GitHub/mytrix/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New(dsn string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.RSSFeed{},
		&model.RSSItem{},
	)
}

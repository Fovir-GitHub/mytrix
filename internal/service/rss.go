package service

import (
	"fmt"

	"codeberg.org/Fovir/mytrix/internal/config"
	"codeberg.org/Fovir/mytrix/internal/repo"
	"gorm.io/gorm"
)

type RSSService interface{}

type RealRSSService struct {
	feedRepo *repo.RSSFeedRepo
	itemRepo *repo.RSSItemRepo
}

func NewRSSService(db *gorm.DB) RSSService {
	cfg := config.Config.RSS
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
	}
}

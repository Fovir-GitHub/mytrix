package service

import (
	"errors"
	"fmt"

	"codeberg.org/Fovir/mytrix/internal/config"
	"codeberg.org/Fovir/mytrix/internal/feed"
	"codeberg.org/Fovir/mytrix/internal/model"
	"codeberg.org/Fovir/mytrix/internal/repo"
	"gorm.io/gorm"
)

type RSSService interface {
	AddFeed(u string) error
	DeleteFeed(id int) error
	Update() ([]model.RSSItem, error)
	AllFeeds() ([]model.RSSFeed, error)
}

type RealRSSService struct {
	feedRepo *repo.RSSFeedRepo
	itemRepo *repo.RSSItemRepo
	parser   *feed.Parser
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
		parser:   feed.New(),
	}
}

func (r *RealRSSService) AddFeed(u string) error {
	feed, _, err := r.parser.ParseURL(u)
	if err != nil {
		return fmt.Errorf("add feed %s failed: %w", u, err)
	}
	return r.feedRepo.Create(feed)
}

func (r *RealRSSService) DeleteFeed(id int) error {
	if err := r.feedRepo.Delete(id); err != nil {
		return fmt.Errorf("delete feed %d failed: %w", id, err)
	}
	return nil
}

func (r *RealRSSService) Update() ([]model.RSSItem, error) {
	var errs []error
	var res []model.RSSItem

	feeds, err := r.AllFeeds()
	if err != nil {
		return nil, fmt.Errorf("get all feeds failed: %w", err)
	}

	for _, feed := range feeds {
		updated, err := r.updateFeed(&feed)
		if err != nil {
			errs = append(errs, err)
		} else {
			res = append(res, updated...)
		}
	}
	return res, errors.Join(errs...)
}

func (r *RealRSSService) updateFeed(feed *model.RSSFeed) ([]model.RSSItem, error) {
	var updated []model.RSSItem

	_, items, err := r.parser.ParseURL(feed.URL)
	if err != nil {
		return nil, fmt.Errorf("parse feed %s failed: %w", feed.URL, err)
	}

	for _, item := range items {
		if err := r.addItem(&item); err == nil {
			updated = append(updated, item)
		}
	}
	return updated, nil
}

func (r *RealRSSService) addItem(item *model.RSSItem) error {
	if err := r.itemRepo.Create(item); err != nil {
		return fmt.Errorf("add item failed: %w", err)
	}
	return nil
}

func (r *RealRSSService) AllFeeds() ([]model.RSSFeed, error) {
	feeds, err := r.feedRepo.AllFeeds()
	if err != nil {
		return nil, fmt.Errorf("fetch all feeds failed: %w", err)
	}
	return feeds, nil
}

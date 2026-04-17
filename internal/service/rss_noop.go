package service

import "codeberg.org/Fovir/mytrix/internal/model"

type NoopRSSService struct {
	err error
}

func (r *NoopRSSService) AddFeed(u string) error {
	return r.err
}

func (r *NoopRSSService) DeleteFeed(id int) error {
	return r.err
}

func (r *NoopRSSService) Update() (string, error) {
	return "", r.err
}

func (r *NoopRSSService) AllFeeds() ([]model.RSSFeed, error) {
	return nil, r.err
}

func (r *NoopRSSService) ListFeeds() (string, error) {
	return "", r.err
}

func (r *NoopRSSService) ExportFeeds() (string, error) {
	return "", r.err
}

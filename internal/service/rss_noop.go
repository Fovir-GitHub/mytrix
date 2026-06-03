package service

import (
	"context"

	"codeberg.org/Fovir/mytrix/internal/db"
)

type NoopRSSService struct {
	err error
}

func (r *NoopRSSService) AddFeeds(context.Context, []string) (string, error) {
	return "", r.err
}

func (r *NoopRSSService) DeleteFeeds(context.Context, []string) (string, error) {
	return "", r.err
}

func (r *NoopRSSService) Update(context.Context) ([]string, error) {
	return nil, r.err
}

func (r *NoopRSSService) AllFeeds(context.Context) ([]db.RssFeed, error) {
	return nil, r.err
}

func (r *NoopRSSService) ListFeeds(context.Context) (string, error) {
	return "", r.err
}

func (r *NoopRSSService) ExportFeeds(context.Context) (string, error) {
	return "", r.err
}

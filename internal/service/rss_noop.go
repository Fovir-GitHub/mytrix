package service

import (
	"context"

	"codeberg.org/Fovir/mytrix/internal/db"
	"codeberg.org/Fovir/mytrix/internal/model"
	"maunium.net/go/mautrix/id"
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

func (r *NoopRSSService) Update(context.Context) ([]model.RSSUpdateResult, error) {
	return nil, r.err
}

func (r *NoopRSSService) AllFeeds(context.Context) ([]db.RSSFeed, error) {
	return nil, r.err
}

func (r *NoopRSSService) ListFeeds(context.Context) (string, error) {
	return "", r.err
}

func (r *NoopRSSService) ExportFeeds(context.Context) (string, error) {
	return "", r.err
}

func (r *NoopRSSService) MarkItemsPushed(context.Context, *model.RSSUpdateResult) error {
	return r.err
}

func (r *NoopRSSService) SetThreadRoot(ctx context.Context, feedID int64, roomID id.RoomID, eventID id.EventID) error {
	return r.err
}

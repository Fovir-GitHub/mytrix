package model

import "codeberg.org/Fovir/mytrix/internal/db"

type RSSUpdateResult struct {
	Rendered string
	ItemIDs  []int64
	Event    *db.Event
	FeedID   int64
}

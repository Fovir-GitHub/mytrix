// Package model contains data models used throughout the application.
package model

import (
	"gorm.io/gorm"
)

// RSSFeed represents an RSS feed subscription.
// It stores the feed URL and title.
type RSSFeed struct {
	gorm.Model
	URL   string `gorm:"uniqueIndex;not null"`
	Title string
}

// RSSItem represents an item from an RSS feed.
// It stores the item's GUID, link, title, and associated feed ID.
type RSSItem struct {
	gorm.Model
	FeedID uint   `gorm:"index:idx_feed_guid,unique"`
	GUID   string `gorm:"index:idx_feed_guid,unique"`
	Link   string `gorm:"index:idx_feed_link,unique"`
	Title  string
}

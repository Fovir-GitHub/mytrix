package model

import (
	"gorm.io/gorm"
)

type RSSFeed struct {
	gorm.Model
	URL   string `gorm:"uniqueIndex;not null"`
	Title string
}

type RSSItem struct {
	gorm.Model
	FeedID uint   `gorm:"index:idx_feed_guid,unique"`
	GUID   string `gorm:"index:idx_feed_guid,unique"`
	Link   string `gorm:"index:idx_feed_link,unique"`
	Title  string
}

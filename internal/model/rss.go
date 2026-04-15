// Package model contains data models used throughout the application.
package model

import (
	"bytes"
	"fmt"

	"gorm.io/gorm"
)

// RSSFeed represents an RSS feed subscription.
// It stores the feed URL and title.
type RSSFeed struct {
	gorm.Model
	URL       string `gorm:"index:idx_feed_url_deleted,unique;not null"`
	Title     string
	DeletedAt gorm.DeletedAt `gorm:"index:idx_feed_url_deleted,unique"`
}

// RSSItem represents an item from an RSS feed.
// It stores the item's GUID, link, title, and associated feed ID.
type RSSItem struct {
	gorm.Model
	FeedID    uint   `gorm:"index"`
	GUID      string `gorm:"index:idx_feed_guid_deleted,unique"`
	Link      string `gorm:"index:idx_feed_link_deleted,unique"`
	Title     string
	DeletedAt gorm.DeletedAt `gorm:"index:idx_feed_guid_deleted,unique;index:idx_feed_link_deleted,unique"`
}

type RSSItemView struct {
	Title string
	Link  string
}

func (r *RSSItem) toView() *RSSItemView {
	return &RSSItemView{
		Title: r.Title,
		Link:  r.Link,
	}
}

func (r *RSSItem) ToMarkdown() string {
	var buf bytes.Buffer
	view := r.toView()
	if err := rssItemTmpl.Execute(&buf, view); err != nil {
		return fmt.Sprintf("Title: %s\nURL: %s", view.Title, view.Link)
	}
	return buf.String()
}

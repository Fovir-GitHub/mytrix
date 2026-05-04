// Package model contains data models used throughout the application.
package model

import (
	"bytes"
	"fmt"
	"strconv"

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
	FeedID      uint   `gorm:"uniqueIndex:idx_feed_guid"`
	GUID        string `gorm:"uniqueIndex:idx_feed_guid"`
	Link        string `gorm:"uniqueIndex"`
	Title       string
	Description string
}

// ToMarkdown returns the RSSFeed formatted as a markdown string using the configured template.
// If template execution fails, it falls back to a simple formatted string representation.
func (r *RSSFeed) ToMarkdown() string {
	var buf bytes.Buffer
	if err := rssFeedTmpl.Execute(&buf, r); err != nil {
		return fmt.Sprintf("ID: %d\nTitle: %s\nURL: %s", r.ID, r.Title, r.URL)
	}
	return buf.String()
}

// ToMarkdown returns the RSSItem formatted as a markdown string using the configured template.
// If template execution fails, it falls back to a simple formatted string representation.
func (r RSSItem) ToMarkdown(feed *RSSFeed) string {
	var buf bytes.Buffer
	titles := []string{
		r.Title,
		feed.Title,
		feed.URL,
		r.Link,
		strconv.Itoa(int(feed.ID)),
	}
	for _, title := range titles {
		if title != "" {
			r.Title = title
			break
		}
	}
	if err := rssItemTmpl.Execute(&buf, r); err != nil {
		return fmt.Sprintf("Title: %s\nURL: %s", r.Title, r.Link)
	}
	return buf.String()
}

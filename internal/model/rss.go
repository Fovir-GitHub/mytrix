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
	URL   string `gorm:"uniqueIndex;not null"`
	Title string
}

// RSSItem represents an item from an RSS feed.
// It stores the item's GUID, link, title, and associated feed ID.
type RSSItem struct {
	gorm.Model
	FeedID uint   `gorm:"uniqueIndex:idx_feed_guid"`
	GUID   string `gorm:"uniqueIndex:idx_feed_guid"`
	Link   string `gorm:"uniqueIndex"`
	Title  string
}

// rssFeedView is a view model for displaying RSS feed information.
type rssFeedView struct {
	ID    uint
	Title string
	URL   string
}

// rssItemView is a view model for displaying RSS item information.
type rssItemView struct {
	Title string
	Link  string
}

func (r *RSSFeed) toView() *rssFeedView {
	return &rssFeedView{
		ID:    r.ID,
		Title: r.Title,
		URL:   r.URL,
	}
}

// ToMarkdown returns the RSSFeed formatted as a markdown string using the configured template.
// If template execution fails, it falls back to a simple formatted string representation.
func (r *RSSFeed) ToMarkdown() string {
	var buf bytes.Buffer
	view := r.toView()
	if err := rssFeedTmpl.Execute(&buf, view); err != nil {
		return fmt.Sprintf("ID: %d\nTitle: %s\nURL: %s", view.ID, view.Title, view.URL)
	}
	return buf.String()
}

func (r *RSSItem) toView() *rssItemView {
	return &rssItemView{
		Title: r.Title,
		Link:  r.Link,
	}
}

// ToMarkdown returns the RSSItem formatted as a markdown string using the configured template.
// If template execution fails, it falls back to a simple formatted string representation.
func (r *RSSItem) ToMarkdown() string {
	var buf bytes.Buffer
	view := r.toView()
	if err := rssItemTmpl.Execute(&buf, view); err != nil {
		return fmt.Sprintf("Title: %s\nURL: %s", view.Title, view.Link)
	}
	return buf.String()
}

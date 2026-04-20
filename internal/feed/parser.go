// Package feed provides functionality for parsing RSS feeds.
// It uses the gofeed library to fetch and parse feed data.
package feed

import (
	"fmt"
	"strings"

	"codeberg.org/Fovir/mytrix/internal/model"
	"codeberg.org/Fovir/mytrix/internal/utils"
	"github.com/mmcdole/gofeed"
)

// Parser wraps gofeed.Parser for fetching and parsing RSS feeds.
type Parser struct {
	p *gofeed.Parser
}

// New returns a new Parser instance initialized with a gofeed.Parser.
// The returned parser is ready to be used for fetching and parsing RSS feeds.
func New() *Parser {
	return &Parser{
		p: gofeed.NewParser(),
	}
}

// ParseURL fetches and parses an RSS feed from the given URL.
// It returns the feed metadata, a list of RSS items, and any error encountered.
// If the URL is unreachable or the feed content is invalid, it returns an error.
func (p *Parser) ParseURL(u string) (*model.RSSFeed, []model.RSSItem, error) {
	feed, err := p.p.ParseURL(u)
	if err != nil {
		return nil, nil, fmt.Errorf("parse url failed (url=%s): %w", u, err)
	}

	rssItems := make([]model.RSSItem, 0, len(feed.Items))
	for _, item := range feed.Items {
		rssItems = append(rssItems, model.RSSItem{
			GUID:        itemGUID(item),
			Link:        item.Link,
			Title:       item.Title,
			Description: item.Description,
		})
	}

	return &model.RSSFeed{
		URL:   feed.FeedLink,
		Title: feed.Title,
	}, rssItems, nil
}

func itemGUID(item *gofeed.Item) string {
	if item.GUID != "" {
		return item.GUID
	}
	if item.Link != "" {
		return item.Link
	}

	var parts strings.Builder
	add := func(s string) {
		if s == "" {
			return
		}
		parts.WriteString(s)
		parts.WriteString("|")
	}

	add(item.Title)
	add(item.Published)
	add(item.Description)
	add(item.Content)

	return utils.HashMD5(parts.String())
}

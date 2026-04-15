// Package feed provides functionality for parsing RSS feeds.
// It uses the gofeed library to fetch and parse feed data.
package feed

import (
	"fmt"

	"codeberg.org/Fovir/mytrix/internal/model"
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
		return nil, nil, fmt.Errorf("parse url %s failed: %w", u, err)
	}

	var rssItems []model.RSSItem
	items := feed.Items
	for _, item := range items {
		rssItem := model.RSSItem{
			GUID:  item.GUID,
			Link:  item.Link,
			Title: item.Title,
		}
		rssItems = append(rssItems, rssItem)
	}

	return &model.RSSFeed{URL: feed.FeedLink, Title: feed.Title}, rssItems, nil
}

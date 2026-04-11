package feed

import (
	"fmt"

	"github.com/Fovir-GitHub/mytrix/internal/model"
	"github.com/mmcdole/gofeed"
)

type Parser struct {
	p *gofeed.Parser
}

func New() *Parser {
	return &Parser{
		p: gofeed.NewParser(),
	}
}

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

	return &model.RSSFeed{URL: feed.Link, Title: feed.Title}, rssItems, nil
}

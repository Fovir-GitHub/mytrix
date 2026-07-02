package render

import (
	"bytes"
	"fmt"
	"log/slog"
	"strconv"

	"codeberg.org/Fovir/mytrix/internal/config"
	"codeberg.org/Fovir/mytrix/internal/db"
)

func RssFeedMarkdown(feed *db.RSSFeed) string {
	var buf bytes.Buffer
	if err := rssFeedTmpl.Execute(&buf, feed); err != nil {
		slog.Debug("rss feed markdown failed", "err", err)
		return fmt.Sprintf("ID: %d\nTitle: %s\nURL: %s", feed.ID, feed.Title, feed.Url)
	}
	return buf.String()
}

func RssItemMarkdown(feed *db.RSSFeed, item *db.RSSItem) string {
	var buf bytes.Buffer
	descMaxLen := config.Config.RSS.DescriptionMaxLength
	titles := []string{
		item.Title,
		feed.Title,
		feed.Url,
		item.Link,
		strconv.Itoa(int(feed.ID)),
	}

	for _, title := range titles {
		if title != "" {
			item.Title = title
			break
		}
	}

	if len(item.Description) > descMaxLen {
		item.Description = ""
	}

	if err := rssItemTmpl.Execute(&buf, item); err != nil {
		return fmt.Sprintf("Title: %s\nURL: %s", item.Title, item.Link)
	}
	return buf.String()
}

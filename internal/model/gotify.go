package model

import (
	"bytes"
	"fmt"
	"log/slog"
	"time"
)

// GotifyMessage represents a notification message from the Gotify service.
// It contains the notification ID, title, message content, and timestamp.
type GotifyMessage struct {
	ID      int       `json:"id"`
	Message string    `json:"message"`
	Title   string    `json:"title"`
	Date    time.Time `json:"date"`
}

type gotifyView struct {
	Title   string
	Message string
	ID      int
	Date    string
}

// ToMarkdown returns the GotifyMessage formatted as a markdown string.
// It formats the message title, content, ID, and date using a template.
// If template execution fails, it falls back to a formatted string representation.
func (g GotifyMessage) ToMarkdown() string {
	var buf bytes.Buffer
	date := g.Date.Format("2006-01-02 15:04:05")
	err := gotifyTmpl.Execute(&buf, gotifyView{
		Title:   g.Title,
		Message: g.Message,
		ID:      g.ID,
		Date:    date,
	})
	if err != nil {
		slog.Error(
			"parse gotify message to markdown failed",
			"title", g.Title,
			"message", g.Message,
			"id", g.ID,
			"date", date,
			"err", err,
		)
		return fmt.Sprintf("Title: %s\n\nMessage: %s\n\nID: %d\nDate: %s", g.Title, g.Message, g.ID, date)
	}
	return buf.String()
}

package render

import (
	"bytes"
	"fmt"
	"log/slog"

	"codeberg.org/Fovir/mytrix/internal/model"
)

type gotifyView struct {
	Title   string
	Message string
	ID      int
	Date    string
}

func GotifyMarkdown(g *model.GotifyMessage) string {
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

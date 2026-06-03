package model

import (
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

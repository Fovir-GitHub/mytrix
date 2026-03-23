package model

import "time"

type GotifyMessage struct {
	Id      int       `json:"id"`
	Message string    `json:"message"`
	Title   string    `json:"title"`
	Date    time.Time `json:"date"`
}

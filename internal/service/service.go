// Package service contains service implementations for various integrations.
package service

type Service struct {
	Gotify  *GotifyService
	Message *MessageService
	Umami   UmamiService
	Wakapi  WakapiService
	RSS     RSSService
}

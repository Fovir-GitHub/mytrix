// Package service contains service implementations for various integrations.
package service

// Service aggregates all service implementations for different integrations like Gotify, Wakapi, Umami, and RSS.
type Service struct {
	Gotify     *GotifyService
	Invitation *InvitationService
	Message    *MessageService
	RSS        RSSService
	Umami      UmamiService
	Wakapi     WakapiService
}

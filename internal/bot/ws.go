// Package bot contains bot-related functionality including WebSocket handling.
package bot

import (
	"net/url"

	"github.com/Fovir-GitHub/mytrix/internal/config"
	"github.com/Fovir-GitHub/mytrix/internal/model"
)

// registerWs registers the Gotify WebSocket with the manager if Gotify is enabled in the configuration.
func (b *Bot) registerWs() {
	b.WsManager.AddIfEnabled(model.SourceGotify, buildGotifyWSURL(), config.Config.Gotify.Enabled)
}

// buildGotifyWSURL builds and returns the WebSocket URL for Gotify.
// It constructs a WSS URL with the Gotify server, /stream path, and token as a query parameter.
func buildGotifyWSURL() string {
	server := config.Config.Gotify.Server
	token := config.Config.Gotify.Token

	u := url.URL{
		Scheme: "wss",
		Host:   server,
		Path:   "/stream",
	}
	q := u.Query()
	q.Set("token", token)
	u.RawQuery = q.Encode()

	return u.String()
}

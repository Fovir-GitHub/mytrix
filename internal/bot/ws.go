package bot

import (
	"net/url"

	"github.com/Fovir-GitHub/mytrix/internal/config"
	"github.com/Fovir-GitHub/mytrix/internal/model"
)

func (b *Bot) registerWs() {
	b.WsManager.AddIfEnabled(model.SourceGotify, buildGotifyWSURL(), config.Config.Gotify.Enabled)
}

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

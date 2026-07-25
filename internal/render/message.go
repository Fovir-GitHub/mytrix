package render

import (
	"codeberg.org/Fovir/mytrix/internal/config"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
)

func RenderMessage(text string) event.MessageEventContent {
	cfg := config.Config.Msg
	return format.RenderMarkdown(text, cfg.AllowMarkdown, cfg.AllowHTML)
}

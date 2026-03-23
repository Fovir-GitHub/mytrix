package handler

import (
	"context"
	"fmt"

	"maunium.net/go/mautrix/event"
)

func (h *Handler) handleGotify(ctx context.Context, evt *event.Event) error {
	const defaultMessage = "no message found"

	msgs, err := h.service.Gotify.FetchMessages()
	if err != nil {
		h.service.Message.Reply(ctx, evt, defaultMessage)
		return fmt.Errorf("fetch message failed: %w", err)
	}

	if len(msgs) == 0 {
		h.service.Message.Reply(ctx, evt, defaultMessage)
	} else {
		h.service.Message.Reply(ctx, evt, msgs[0].Message)
	}

	return nil
}

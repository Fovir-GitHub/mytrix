package matrix

import (
	"context"
	"fmt"

	"github.com/Fovir-GitHub/mytrix/internal/config"
	"github.com/Fovir-GitHub/mytrix/internal/crypto"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto/cryptohelper"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
	"maunium.net/go/mautrix/id"
)

type Client struct {
	c *mautrix.Client
}

func New(c *mautrix.Client) *Client {
	return &Client{c: c}
}

func (m *Client) SendTextMessage(ctx context.Context, roomID id.RoomID, text string) error {
	cfg := config.Config.Msg
	content := format.RenderMarkdown(text, cfg.AllowMarkdown, cfg.AllowHTML)
	_, err := m.c.SendMessageEvent(ctx, roomID, event.EventMessage, content)
	return err
}

func (m *Client) Sync() error {
	return m.c.Sync()
}

func (m *Client) VerifyWithRecoveryKey() error {
	ch, ok := m.c.Crypto.(*cryptohelper.CryptoHelper)
	if !ok {
		return fmt.Errorf("crypto helper type mismatch")
	}

	return crypto.VerifyWithRecoveryKey(ch.Machine())
}

func (m *Client) UserID() id.UserID {
	return m.c.UserID
}

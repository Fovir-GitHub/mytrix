package matrix

import (
	"context"
	"fmt"

	"codeberg.org/Fovir/mytrix/internal/config"
	"codeberg.org/Fovir/mytrix/internal/crypto"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto/cryptohelper"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
	"maunium.net/go/mautrix/id"
)

// Client is a wrapper around mautrix.Client that provides Matrix protocol operations.
type Client struct {
	c *mautrix.Client
}

// New creates a new Client wrapping the provided mautrix client.
func New(c *mautrix.Client) *Client {
	return &Client{c: c}
}

// SendTextMessage sends a formatted text message to the specified room.
// It converts Markdown and HTML according to the message configuration.
func (m *Client) SendTextMessage(ctx context.Context, roomID id.RoomID, text string) error {
	cfg := config.Config.Msg
	content := format.RenderMarkdown(text, cfg.AllowMarkdown, cfg.AllowHTML)
	_, err := m.c.SendMessageEvent(ctx, roomID, event.EventMessage, content)
	return err
}

// Sync performs a full synchronization with the Matrix homeserver.
func (m *Client) Sync() error {
	return m.c.Sync()
}

// VerifyWithRecoveryKey verifies the end-to-end encryption setup using the recovery key.
func (m *Client) VerifyWithRecoveryKey() error {
	ch, ok := m.c.Crypto.(*cryptohelper.CryptoHelper)
	if !ok {
		return fmt.Errorf("crypto helper type mismatch")
	}

	return crypto.VerifyWithRecoveryKey(ch.Machine())
}

// UserID returns the user ID of the Matrix client.
func (m *Client) UserID() id.UserID {
	return m.c.UserID
}

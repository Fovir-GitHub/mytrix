package bot

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Fovir-GitHub/mytrix/internal/config"
	"github.com/Fovir-GitHub/mytrix/internal/storage"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"
)

// newClient creates a new Matrix client, either by loading an existing session
// or by logging in with credentials if no session exists.
// It returns the initialized client and any error encountered.
func newClient() (*mautrix.Client, error) {
	cfg := config.Config
	var client *mautrix.Client

	session, err := storage.LoadSession()
	if err == nil {
		slog.Info("login with existed session")
		client, err = mautrix.NewClient(cfg.Homeserver, id.UserID(session.UserID), session.AccessToken)
		if err != nil {
			return nil, fmt.Errorf("create mautrix client failed: %w", err)
		}
		client.DeviceID = id.DeviceID(session.DeviceID)
	} else {
		slog.Info(
			"session does not exist. use login",
			"err", err,
		)
		client, err = mautrix.NewClient(cfg.Homeserver, "", "")
		if err != nil {
			return nil, fmt.Errorf("create mautrix client failed: %w", err)
		}
		resp, err := login(client)
		if err != nil {
			return nil, fmt.Errorf("login failed: %w", err)
		}
		client.UserID = resp.UserID
		client.AccessToken = resp.AccessToken
		client.DeviceID = resp.DeviceID

		_ = storage.SaveSession(resp)
	}

	return client, nil
}

// login performs a password-based login to the Matrix homeserver.
// It takes a client instance and logs in using the bot's credentials from config.
// It returns the login response and any error encountered.
func login(client *mautrix.Client) (*mautrix.RespLogin, error) {
	cfg := config.Config

	resp, err := client.Login(context.Background(), &mautrix.ReqLogin{
		Type: mautrix.AuthTypePassword,
		Identifier: mautrix.UserIdentifier{
			User: cfg.Bot.Username,
			Type: mautrix.IdentifierTypeUser,
		},
		Password:         cfg.Bot.Password,
		StoreCredentials: true,
	})
	if err != nil {
		return nil, fmt.Errorf("login failed: %w", err)
	}

	return resp, nil
}

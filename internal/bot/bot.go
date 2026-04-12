// Package bot contains bot-related functionality.
package bot

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Fovir-GitHub/mytrix/internal/crypto"
	"github.com/Fovir-GitHub/mytrix/internal/handler"
	myhttp "github.com/Fovir-GitHub/mytrix/internal/http"
	"github.com/Fovir-GitHub/mytrix/internal/matrix"
	"github.com/Fovir-GitHub/mytrix/internal/scheduler"
	"github.com/Fovir-GitHub/mytrix/internal/service"
	"github.com/Fovir-GitHub/mytrix/internal/ws"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

// Bot represents a Matrix bot client with sync and encryption support.
type Bot struct {
	Client    *matrix.Client
	WsManager *ws.Manager
	Syncer    *mautrix.DefaultSyncer
	Ready     chan struct{}
	Handler   *handler.Handler
	Scheduler *scheduler.Scheduler
}

// New creates and initializes a Bot instance.
// It sets up the Matrix client, syncer, and encryption helper.
// After creation, call Start() to begin syncing.
func New() (*Bot, error) {
	slog.Debug("new bot start")

	client, err := newClient()
	if err != nil {
		return nil, fmt.Errorf("create client failed: %w", err)
	}

	syncer, ready := setupSyncer()
	client.Syncer = syncer

	cryptoHelper, err := crypto.SetupCryptoHelper(client)
	if err != nil {
		return nil, fmt.Errorf("create cryptohelper failed: %w", err)
	}
	client.Crypto = cryptoHelper

	matrixClient := matrix.New(client)
	http := myhttp.New()
	scheduler := scheduler.NewScheduler()
	service := service.NewService(http, matrixClient, scheduler)
	handler := handler.NewHandler(service)
	wsManager := ws.NewManager()

	bot := &Bot{
		Client:    matrixClient,
		WsManager: wsManager,
		Syncer:    syncer,
		Ready:     ready,
		Handler:   handler,
		Scheduler: scheduler,
	}

	bot.registerHandler()
	bot.registerWs()
	bot.registerScheduler()

	return bot, nil
}

// Start begins the bot's operation.
// It starts the scheduler, Matrix sync, and WebSocket event handling.
// The function blocks until the context is cancelled.
func (b *Bot) Start(ctx context.Context) error {
	b.Scheduler.Start()

	go func() {
		if err := b.Client.Sync(); err != nil {
			slog.Error(
				"runtime error",
				"err", err,
			)
			return
		}
	}()

	go func() {
		for event := range b.WsManager.Events() {
			slog.Debug(
				"receive websocket event",
				"source", event.Source,
				"data", string(event.Data),
			)
			err := b.Handler.HandleWSEvent(ctx, event)
			if err != nil {
				slog.Error(
					"handle websocket event failed",
					"err", err,
				)
			}
		}
	}()

	<-b.Ready

	err := b.Client.VerifyWithRecoveryKey()
	if err != nil {
		return fmt.Errorf("verify recovery key failed: %w", err)
	}

	<-ctx.Done()
	return ctx.Err()
}

// registerHandler registers the command handler for Matrix message events by setting up
// the syncer to call HandleCommand when a message event is received.
func (b *Bot) registerHandler() {
	b.Syncer.OnEventType(event.EventMessage, b.Handler.HandleCommand)
}

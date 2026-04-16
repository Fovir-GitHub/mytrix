// Package bot contains bot-related functionality.
package bot

import (
	"context"
	"fmt"
	"log/slog"

	"codeberg.org/Fovir/mytrix/internal/crypto"
	"codeberg.org/Fovir/mytrix/internal/handler"
	myhttp "codeberg.org/Fovir/mytrix/internal/http"
	"codeberg.org/Fovir/mytrix/internal/matrix"
	"codeberg.org/Fovir/mytrix/internal/scheduler"
	"codeberg.org/Fovir/mytrix/internal/service"
	"codeberg.org/Fovir/mytrix/internal/ws"
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
	client, err := newClient()
	if err != nil {
		return nil, fmt.Errorf("create bot failed: %w", err)
	}

	syncer, ready := setupSyncer()
	client.Syncer = syncer

	cryptoHelper, err := crypto.SetupCryptoHelper(client)
	if err != nil {
		return nil, fmt.Errorf("create bot failed: %w", err)
	}
	client.Crypto = cryptoHelper

	db, err := setupDB()
	if err != nil {
		return nil, fmt.Errorf("create bot failed: %w", err)
	}

	matrixClient := matrix.New(client)
	http := myhttp.New()
	scheduler := scheduler.NewScheduler()

	gotifySrv := service.NewGotifyService()
	msgSrv := service.NewMessageService(matrixClient)
	umamiSrv := service.NewUmamiService(http)
	wakapiSrv := service.NewWakapiService(http, scheduler)
	rssSrv := service.NewRSSService(db)

	service := &service.Service{
		Gotify:  gotifySrv,
		Message: msgSrv,
		Umami:   umamiSrv,
		Wakapi:  wakapiSrv,
		RSS:     rssSrv,
	}

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
				"source", event.Source)
			err := b.Handler.HandleWSEvent(ctx, event)
			if err != nil {
				slog.Error(
					"handle websocket event failed",
					"source", event.Source,
					"err", err,
				)
			}
		}
	}()

	<-b.Ready

	err := b.Client.VerifyWithRecoveryKey()
	if err != nil {
		return fmt.Errorf("start bot failed: %w", err)
	}

	<-ctx.Done()
	return ctx.Err()
}

// registerHandler registers the command handler for Matrix message events by setting up
// the syncer to call HandleCommand when a message event is received.
func (b *Bot) registerHandler() {
	b.Syncer.OnEventType(event.EventMessage, b.Handler.HandleCommand)
}

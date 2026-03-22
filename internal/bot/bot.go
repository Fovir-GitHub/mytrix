package bot

import (
	"fmt"
	"log/slog"

	clientpkg "github.com/Fovir-GitHub/mytrix/internal/client"
	"github.com/Fovir-GitHub/mytrix/internal/crypto"
	"github.com/Fovir-GitHub/mytrix/internal/handler"
	"github.com/Fovir-GitHub/mytrix/internal/service"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto/cryptohelper"
)

// Bot represents a Matrix bot client with sync and encryption support.
type Bot struct {
	Client  *clientpkg.MatrixClient
	Syncer  *mautrix.DefaultSyncer
	Ready   chan struct{}
	Handler *handler.MessageHandler
}

// New creates and initializes a `Bot` instance.
// It sets up the Matrix client, syncer, and encryption helper.
// After creation, call `Start()` to begin syncing.
func New() (*Bot, error) {
	slog.Debug("new bot start")

	client, err := newClient()
	if err != nil {
		return nil, fmt.Errorf("create client failed: %w", err)
	}

	syncer, ready := setupSyncer()
	client.Syncer = syncer

	cryptoHelper, err := setupCryptoHelper(client)
	if err != nil {
		return nil, fmt.Errorf("create cryptohelper failed: %w", err)
	}
	client.Crypto = cryptoHelper

	matrixClient := clientpkg.New(client)
	messageService := service.NewMessageService(matrixClient)
	messageHandler := handler.NewMessageHandler(messageService)

	bot := &Bot{
		Client:  matrixClient,
		Syncer:  syncer,
		Ready:   ready,
		Handler: messageHandler,
	}

	bot.registerHandler()

	return bot, nil
}

func (b *Bot) Start() error {
	go func() {
		if err := b.Client.Sync(); err != nil {
			panic(err)
		}
	}()

	<-b.Ready

	err := verifyWithRecoveryKey(b.Client.Crypto.(*cryptohelper.CryptoHelper).Machine())
	if err != nil {
		return fmt.Errorf("verify recovery key failed: %w", err)
	}
	return nil
}

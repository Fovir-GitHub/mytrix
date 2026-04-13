// Package crypto provides functions for setting up cryptographic support in the bot.
// It handles encryption setup for Matrix client sessions.
package crypto

import (
	"context"
	"fmt"
	"path/filepath"

	"codeberg.org/Fovir/mytrix/internal/config"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto"
	"maunium.net/go/mautrix/crypto/cryptohelper"
)

// SetupCryptoHelper sets up and returns a new CryptoHelper for the given Matrix client.
// It initializes encryption support using the configured pickle key and stores the database
// in the bot's data directory. The caller is responsible for managing the helper's lifecycle.
func SetupCryptoHelper(client *mautrix.Client) (*cryptohelper.CryptoHelper, error) {
	dbPath := filepath.Join(config.Config.Datadir, "crypto.db")
	helper, err := cryptohelper.NewCryptoHelper(client, []byte(config.Config.Bot.PickleKey), dbPath)
	if err != nil {
		return nil, fmt.Errorf("create cryptohelper failed: %w", err)
	}
	err = helper.Init(context.Background())
	if err != nil {
		return nil, fmt.Errorf("init cryptohelper failed: %w", err)
	}
	return helper, nil
}

// verifyWithRecoveryKey verifies the recovery key and completes the encryption setup.
// It fetches cross-signing keys from the SSSS service and signs the device and master key.
func VerifyWithRecoveryKey(machine *crypto.OlmMachine) (err error) {
	ctx := context.Background()
	keyID, keyData, err := machine.SSSS.GetDefaultKeyData(ctx)
	if err != nil {
		return err
	}
	key, err := keyData.VerifyRecoveryKey(keyID, config.Config.Bot.RecoveryKey)
	if err != nil {
		return err
	}
	err = machine.FetchCrossSigningKeysFromSSSS(ctx, key)
	if err != nil {
		return err
	}
	err = machine.SignOwnDevice(ctx, machine.OwnIdentity())
	if err != nil {
		return err
	}
	err = machine.SignOwnMasterKey(ctx)
	return err
}

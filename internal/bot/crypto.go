package bot

import (
	"context"
	"fmt"

	"github.com/Fovir-GitHub/mytrix/internal/config"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto"
	"maunium.net/go/mautrix/crypto/cryptohelper"
)

// setupCryptoHelper creates and initializes a new CryptoHelper for the Matrix client.
// It sets up encryption support using the provided pickle key and database path.
func setupCryptoHelper(client *mautrix.Client) (*cryptohelper.CryptoHelper, error) {
	helper, err := cryptohelper.NewCryptoHelper(client, []byte(config.Config.Bot.PickleKey), "db/crypto.db")
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
func verifyWithRecoveryKey(machine *crypto.OlmMachine) (err error) {
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

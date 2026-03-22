package bot

import (
	"context"
	"sync"

	"maunium.net/go/mautrix"
)

// setupSyncer creates and configures a Matrix client syncer with a ready channel.
// It returns a configured syncer and a channel that gets closed when the initial sync completes.
func setupSyncer() (*mautrix.DefaultSyncer, chan struct{}) {
	readyChan := make(chan struct{})
	var once sync.Once

	syncer := mautrix.NewDefaultSyncer()
	syncer.OnSync(func(ctx context.Context, resp *mautrix.RespSync, since string) bool {
		once.Do(func() {
			close(readyChan)
		})
		return true
	})

	return syncer, readyChan
}

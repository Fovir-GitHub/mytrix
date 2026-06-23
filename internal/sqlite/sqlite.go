package sqlite

import (
	"database/sql"
	"fmt"

	"codeberg.org/Fovir/mytrix/internal/config"
)

func Open(dsn string) (*sql.DB, error) {
	cfg := config.Config.DB
	pragmas := fmt.Sprintf("?_journal_mode=WAL&_busy_timeout=%d", cfg.BusyTimeout)
	return sql.Open("sqlite", dsn+pragmas)
}

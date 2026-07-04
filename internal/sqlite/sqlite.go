package sqlite

import (
	"database/sql"
	"fmt"

	"codeberg.org/Fovir/mytrix/internal/config"
)

func Open(dsn string) (*sql.DB, error) {
	cfg := config.Config.DB
	pragmas := fmt.Sprintf(
		"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(%d)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(ON)",
		cfg.BusyTimeout,
	)
	db, err := sql.Open("sqlite", dsn+pragmas)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	return db, nil
}

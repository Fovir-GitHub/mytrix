package sqlite

import "database/sql"

const pragmas = "?_journal_mode=WAL&_busy_timeout=5000"

func Open(dsn string) (*sql.DB, error) {
	return sql.Open("sqlite", dsn+pragmas)
}

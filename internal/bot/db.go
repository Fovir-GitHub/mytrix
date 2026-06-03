package bot

import (
	"context"
	"database/sql"
	"fmt"
	"path"

	_ "embed"

	"codeberg.org/Fovir/mytrix/internal/assets"
	"codeberg.org/Fovir/mytrix/internal/config"
	"codeberg.org/Fovir/mytrix/internal/db"
)

// setupDB initializes the database connection and runs migrations.
// It constructs the DSN from the configured data directory and database path.
func setupDB() (*db.Queries, error) {
	cfg := config.Config
	dsn := path.Join(cfg.Datadir, cfg.DBPath)
	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("open database failed (dsn=%s): %w", dsn, err)
	}

	entries, err := assets.SchemaFS.ReadDir(assets.SchemaPath)
	if err != nil {
		return nil, fmt.Errorf("read schema dir failed (path=%s): %w", assets.SchemaPath, err)
	}
	for _, entry := range entries {
		content, _ := assets.SchemaFS.ReadFile(path.Join(assets.SchemaPath, entry.Name()))
		if _, err := conn.ExecContext(context.Background(), string(content)); err != nil {
			return nil, err
		}
	}

	q := db.New(conn)
	return q, nil
}

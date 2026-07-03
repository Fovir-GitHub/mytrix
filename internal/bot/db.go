package bot

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"path"

	_ "embed"

	"codeberg.org/Fovir/mytrix/internal/assets"
	"codeberg.org/Fovir/mytrix/internal/config"
	"codeberg.org/Fovir/mytrix/internal/db"
	"codeberg.org/Fovir/mytrix/internal/sqlite"
)

// setupDB initializes the database connection and runs migrations.
// It constructs the DSN from the configured data directory and database path.
func setupDB() (sqlDB *sql.DB, q *db.Queries, err error) {
	cfg := config.Config
	dsn := path.Join(cfg.Datadir, cfg.DB.Path)

	sqlDB, err = sqlite.Open(dsn)
	if err != nil {
		_ = sqlDB.Close()
		return nil, nil, fmt.Errorf("open database failed (dsn=%s): %w", dsn, err)
	}

	entries, err := assets.SchemaFS.ReadDir(assets.SchemaPath)
	if err != nil {
		return nil, nil, fmt.Errorf("read schema dir failed (path=%s): %w", assets.SchemaPath, err)
	}

	for _, entry := range entries {
		content, _ := assets.SchemaFS.ReadFile(path.Join(assets.SchemaPath, entry.Name()))
		if _, err := sqlDB.ExecContext(context.Background(), string(content)); err != nil {
			return nil, nil, err
		}
	}

	if err := runMigrations(sqlDB); err != nil {
		return nil, nil, err
	}
	slog.Info("database migration finished")

	q = db.New(sqlDB)
	return sqlDB, q, nil
}

func syncConfig(ctx context.Context, q *db.Queries) error {
	cfg := config.Config

	// Add admin information into the database.
	if err := q.CreateUser(ctx, &db.CreateUserParams{
		ID:   cfg.User.AdminID,
		Role: "admin",
	}); err != nil {
		return err
	}

	if err := q.CreateRoom(ctx, &db.CreateRoomParams{
		ID:    cfg.RoomID,
		State: "joined",
	}); err != nil {
		return err
	}

	return nil
}

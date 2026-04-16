package bot

import (
	"fmt"
	"path"

	"codeberg.org/Fovir/mytrix/internal/config"
	"codeberg.org/Fovir/mytrix/internal/database"
	"gorm.io/gorm"
)

func setupDB() (*gorm.DB, error) {
	cfg := config.Config
	dsn := path.Join(cfg.Datadir, cfg.DBPath)
	db, err := database.New(dsn)
	if err != nil {
		return nil, fmt.Errorf("setup database failed (dsn=%s): %w", dsn, err)
	}
	if err := database.Migrate(db); err != nil {
		return nil, fmt.Errorf("setup database failed: %w", err)
	}
	return db, nil
}

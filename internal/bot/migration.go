package bot

import (
	"database/sql"
	"fmt"
	"log/slog"
	"path"

	"codeberg.org/Fovir/mytrix/internal/assets"
)

// runMigrations must run after creating schema.
func runMigrations(db *sql.DB) error {
	entries, err := assets.MigrationsFS.ReadDir(assets.MigrationsPath)
	if err != nil {
		return fmt.Errorf("read migration dir failed (path=%s): %w", assets.MigrationsPath, err)
	}

	rows, err := db.Query(`SELECT name FROM migration`)
	if err != nil {
		return fmt.Errorf("query completed migrations failed: %w", err)
	}
	defer rows.Close() //nolint

	completed := make(map[string]struct{})

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return err
		}
		completed[name] = struct{}{}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	for _, entry := range entries {
		filename := entry.Name()
		if _, ok := completed[filename]; ok {
			continue
		}

		slog.Debug("run migration", "filename", filename)

		sqlPath := path.Join(assets.MigrationsPath, filename)
		migrationSql, err := assets.MigrationsFS.ReadFile(sqlPath)
		if err != nil {
			return fmt.Errorf("read migration SQL file failed (path=%s): %w", sqlPath, err)
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("create transaction failed: %w", err)
		}

		_, err = tx.Exec(string(migrationSql))
		if err != nil {
			if errTx := tx.Rollback(); errTx != nil {
				return fmt.Errorf("rollback %v failed: %w", filename, errTx)
			}
			return fmt.Errorf("exec migration %v failed: %w", filename, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %v failed: %w", filename, err)
		}

		if _, err := db.Exec(`INSERT OR IGNORE INTO migration (name) VALUES (?)`, filename); err != nil {
			return fmt.Errorf("update migration table failed: %w", err)
		}
	}

	return nil
}

// Package database provides database initialization and migration functions.
// It handles SQLite database connections and schema migrations.
package database

import (
	"codeberg.org/Fovir/mytrix/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// New returns a new GORM database connection using the SQLite driver.
// It opens a connection to the database specified by the DSN and returns the connection
// along with any error encountered.
func New(dsn string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
}

// Migrate runs database migrations for the application models.
// It automatically migrates the RSSFeed and RSSItem schemas using GORM's AutoMigrate function
// and returns any error encountered during the migration process.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.RSSFeed{},
		&model.RSSItem{},
	)
}

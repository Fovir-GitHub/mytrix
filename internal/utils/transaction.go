package utils

import (
	"context"
	"database/sql"
	"fmt"

	"codeberg.org/Fovir/mytrix/internal/db"
)

func CreateTransaction(ctx context.Context, sqlDB *sql.DB, q *db.Queries, opts *sql.TxOptions) (*sql.Tx, *db.Queries, error) {
	tx, err := sqlDB.BeginTx(ctx, opts)
	if err != nil {
		return nil, nil, fmt.Errorf("create tx failed: %w", err)
	}
	qtx := q.WithTx(tx)
	return tx, qtx, nil
}

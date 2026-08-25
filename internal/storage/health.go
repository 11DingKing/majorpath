package storage

import (
	"context"
	"database/sql"
	"fmt"
)

func Ready(ctx context.Context, db *sql.DB) error {
	if e := db.PingContext(ctx); e != nil {
		return fmt.Errorf("database ping: %w", e)
	}
	var n int
	if e := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations").Scan(&n); e != nil {
		return fmt.Errorf("migration check: %w", e)
	}
	if n == 0 {
		return fmt.Errorf("no migrations applied")
	}
	return nil
}

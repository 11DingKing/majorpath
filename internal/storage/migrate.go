package storage

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func Migrate(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS schema_migrations(version INTEGER PRIMARY KEY, applied_at TEXT NOT NULL)"); err != nil {
		return fmt.Errorf("create migration table: %w", err)
	}
	entries, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	for _, e := range entries {
		b, er := migrationFiles.ReadFile("migrations/" + e.Name())
		if er != nil {
			return er
		}
		parts := strings.SplitN(e.Name(), "_", 2)
		v, _ := strconv.Atoi(parts[0])
		var applied int
		if er = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version=?", v).Scan(&applied); er != nil {
			return er
		}
		if applied > 0 {
			continue
		}
		tx, er := db.BeginTx(ctx, nil)
		if er != nil {
			return er
		}
		if _, er = tx.ExecContext(ctx, string(b)); er != nil {
			tx.Rollback()
			return fmt.Errorf("migration %s: %w", e.Name(), er)
		}
		if _, er = tx.ExecContext(ctx, "INSERT INTO schema_migrations(version, applied_at) VALUES(?,datetime('now'))", v); er != nil {
			tx.Rollback()
			return er
		}
		if er = tx.Commit(); er != nil {
			return er
		}
	}
	return nil
}

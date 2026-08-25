package repository

import (
	"context"
	"database/sql"
	"time"
)

type CleanupRepo struct{ DB *sql.DB }

func (r CleanupRepo) PurgeExpiredSessions(ctx context.Context, before time.Time) (int64, error) {
	res, e := r.DB.ExecContext(ctx, "DELETE FROM sessions WHERE revoked_at IS NOT NULL OR expires_at<?", before.UTC().Format("2006-01-02 15:04:05"))
	if e != nil {
		return 0, wrap("purge sessions", e)
	}
	return res.RowsAffected()
}
func (r CleanupRepo) PurgeIdempotency(ctx context.Context, now time.Time) (int64, error) {
	res, e := r.DB.ExecContext(ctx, "DELETE FROM idempotency_keys WHERE expires_at<?", now.UTC().Format("2006-01-02 15:04:05"))
	if e != nil {
		return 0, wrap("purge idempotency", e)
	}
	return res.RowsAffected()
}
func (r CleanupRepo) Vacuum(ctx context.Context) error {
	_, e := r.DB.ExecContext(ctx, "PRAGMA optimize")
	return e
}

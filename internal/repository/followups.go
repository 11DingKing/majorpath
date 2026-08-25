package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/majorpath/internal/domain"
	"time"
)

type FollowupRepo struct{ DB *sql.DB }

func (r FollowupRepo) Create(ctx context.Context, id, rec string, due time.Time) error {
	_, e := r.DB.ExecContext(ctx, "INSERT INTO followups(id,recommendation_id,due_at,status,created_at) VALUES(?,?,?,'pending',datetime('now'))", id, rec, due.UTC().Format(time.RFC3339))
	return wrap("followup", e)
}
func (r FollowupRepo) Claim(ctx context.Context, now, lease time.Time) (domain.Followup, error) {
	tx, e := r.DB.BeginTx(ctx, nil)
	if e != nil {
		return domain.Followup{}, e
	}
	defer tx.Rollback()
	var f domain.Followup
	row := tx.QueryRowContext(ctx, "SELECT id,recommendation_id,due_at,status,attempts FROM followups WHERE status IN ('pending','failed') AND due_at<=? AND (lease_until IS NULL OR lease_until<?) ORDER BY due_at LIMIT 1", now.UTC().Format(time.RFC3339), now.UTC().Format(time.RFC3339))
	if e = row.Scan(&f.ID, &f.RecommendationID, &f.DueAt, &f.Status, &f.Attempts); e != nil {
		return f, wrap("claim", e)
	}
	if _, e = tx.ExecContext(ctx, "UPDATE followups SET status='running',lease_until=?,attempts=attempts+1 WHERE id=?", lease.UTC().Format(time.RFC3339), f.ID); e != nil {
		return f, e
	}
	if e = tx.Commit(); e != nil {
		return f, e
	}
	f.Status = domain.FollowupRunning
	f.Attempts++
	return f, nil
}
func (r FollowupRepo) Complete(ctx context.Context, id string, ok bool, msg string) error {
	status := "done"
	if !ok {
		status = "failed"
	}
	_, e := r.DB.ExecContext(ctx, "UPDATE followups SET status=?,last_error=?,completed_at=CASE WHEN ?='done' THEN datetime('now') ELSE completed_at END,lease_until=NULL WHERE id=?", status, msg, status, id)
	return wrap("complete followup", e)
}

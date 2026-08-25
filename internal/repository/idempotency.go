package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/majorpath/internal/domain"
)

type IdempotencyRepo struct{ DB *sql.DB }

func (r IdempotencyRepo) Get(ctx context.Context, key, actor string) (string, error) {
	var body string
	e := r.DB.QueryRowContext(ctx, "SELECT response_body FROM idempotency_keys WHERE key=? AND actor_id=? AND expires_at>datetime('now')", key, actor).Scan(&body)
	return body, wrap("idempotency", e)
}
func (r IdempotencyRepo) Put(ctx context.Context, key, actor, hash, body string) error {
	_, e := r.DB.ExecContext(ctx, "INSERT INTO idempotency_keys(key,actor_id,request_hash,response_body,created_at,expires_at) VALUES(?,?,?,?,datetime('now'),datetime('now','+1 day'))", key, actor, hash, body)
	return wrap("idempotency", e)
}

var _ = domain.ErrConflict

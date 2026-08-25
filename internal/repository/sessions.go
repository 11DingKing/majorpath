package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/majorpath/internal/domain"
)

type SessionRepo struct{ DB *sql.DB }

func (r SessionRepo) Create(ctx context.Context, id, user, hash string, expires string) error {
	_, e := r.DB.ExecContext(ctx, "INSERT INTO sessions(id,user_id,token_hash,expires_at,created_at) VALUES(?,?,?,?,datetime('now'))", id, user, hash, expires)
	return wrap("create session", e)
}
func (r SessionRepo) Find(ctx context.Context, hash string) (string, error) {
	var id string
	e := r.DB.QueryRowContext(ctx, "SELECT user_id FROM sessions WHERE token_hash=? AND revoked_at IS NULL AND expires_at>datetime('now')", hash).Scan(&id)
	return id, wrap("find session", e)
}
func (r SessionRepo) Revoke(ctx context.Context, hash string) error {
	_, e := r.DB.ExecContext(ctx, "UPDATE sessions SET revoked_at=datetime('now') WHERE token_hash=? AND revoked_at IS NULL", hash)
	return wrap("revoke session", e)
}
func (r SessionRepo) Expire(ctx context.Context) error {
	_, e := r.DB.ExecContext(ctx, "UPDATE sessions SET revoked_at=datetime('now') WHERE revoked_at IS NULL AND expires_at<=datetime('now')")
	return wrap("expire sessions", e)
}

var _ = domain.ErrExpired

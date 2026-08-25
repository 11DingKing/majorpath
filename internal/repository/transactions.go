package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type TxRunner struct{ DB *sql.DB }

func (r TxRunner) Run(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, e := r.DB.BeginTx(ctx, nil)
	if e != nil {
		return fmt.Errorf("begin: %w", e)
	}
	if e = fn(tx); e != nil {
		_ = tx.Rollback()
		return e
	}
	if e = tx.Commit(); e != nil {
		return fmt.Errorf("commit: %w", e)
	}
	return nil
}
func ExecTx(ctx context.Context, tx *sql.Tx, query string, args ...any) error {
	_, e := tx.ExecContext(ctx, query, args...)
	return e
}

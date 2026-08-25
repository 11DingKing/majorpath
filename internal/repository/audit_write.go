package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type AuditWriter struct{ DB *sql.DB }

func (r AuditWriter) Write(ctx context.Context, actor, entity, id, action, result, request, detail string) error {
	if entity == "" || action == "" {
		return fmt.Errorf("audit fields required")
	}
	_, e := r.DB.ExecContext(ctx, "INSERT INTO audit_events(id,actor_id,entity_type,entity_id,action,result,request_id,detail,created_at) VALUES(lower(hex(randomblob(16))),?,?,?,?,?,?,?,datetime('now'))", actor, entity, id, action, result, request, detail)
	return wrap("audit write", e)
}
func (r AuditWriter) Count(ctx context.Context, entity string) (int, error) {
	var n int
	e := r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM audit_events WHERE entity_type=?", entity).Scan(&n)
	return n, wrap("audit count", e)
}

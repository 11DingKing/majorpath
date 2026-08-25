package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/majorpath/internal/domain"
)

type AuditRepo struct{ DB *sql.DB }

func (r AuditRepo) List(ctx context.Context, entity string) ([]domain.AuditEvent, error) {
	rows, e := r.DB.QueryContext(ctx, "SELECT id,COALESCE(actor_id,''),entity_type,entity_id,action,result,request_id,detail,created_at FROM audit_events WHERE entity_type=? ORDER BY created_at DESC LIMIT 100", entity)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []domain.AuditEvent{}
	for rows.Next() {
		var a domain.AuditEvent
		if e = rows.Scan(&a.ID, &a.ActorID, &a.EntityType, &a.EntityID, &a.Action, &a.Result, &a.RequestID, &a.Detail, &a.CreatedAt); e != nil {
			return nil, e
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

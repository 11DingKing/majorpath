package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/11DingKing/majorpath/internal/domain"
)

type ReportRepo struct{ DB *sql.DB }
type RecommendationSummary struct {
	ID, Student, Status string
	Items               int
	UpdatedAt           string
}

func (r ReportRepo) Summaries(ctx context.Context, owner string, limit, offset int) ([]RecommendationSummary, error) {
	if limit < 1 || limit > 100 {
		limit = 20
	}
	q := `SELECT r.id,s.name,r.status,COUNT(i.id),r.updated_at FROM recommendations r JOIN students s ON s.id=r.student_id LEFT JOIN recommendation_items i ON i.recommendation_id=r.id WHERE s.owner_id=? GROUP BY r.id ORDER BY r.updated_at DESC LIMIT ? OFFSET ?`
	rows, e := r.DB.QueryContext(ctx, q, owner, limit, offset)
	if e != nil {
		return nil, wrap("summary", e)
	}
	defer rows.Close()
	out := []RecommendationSummary{}
	for rows.Next() {
		var x RecommendationSummary
		if e = rows.Scan(&x.ID, &x.Student, &x.Status, &x.Items, &x.UpdatedAt); e != nil {
			return nil, e
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r ReportRepo) ReservePlan(ctx context.Context, id string) error {
	res, e := r.DB.ExecContext(ctx, "UPDATE admission_plans SET remaining=remaining-1 WHERE id=? AND remaining>0", id)
	if e != nil {
		return wrap("reserve plan", e)
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return fmt.Errorf("%w: no quota", domain.ErrConflict)
	}
	return nil
}
func (r ReportRepo) ReleasePlan(ctx context.Context, id string) error {
	res, e := r.DB.ExecContext(ctx, "UPDATE admission_plans SET remaining=remaining+1 WHERE id=? AND remaining<quota", id)
	if e != nil {
		return wrap("release plan", e)
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return fmt.Errorf("%w: invalid release", domain.ErrConflict)
	}
	return nil
}

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/11DingKing/majorpath/internal/domain"
)

type RecommendationRepo struct{ DB *sql.DB }
type RecommendationInput struct {
	ID, StudentID, CreatedBy, Note string
	Items                          []domain.RecommendationItem
}

func (r RecommendationRepo) Create(ctx context.Context, in RecommendationInput, requestID string) (domain.Recommendation, error) {
	tx, e := r.DB.BeginTx(ctx, nil)
	if e != nil {
		return domain.Recommendation{}, e
	}
	defer tx.Rollback()
	var owner string
	if e = tx.QueryRowContext(ctx, "SELECT owner_id FROM students WHERE id=?", in.StudentID).Scan(&owner); e != nil {
		return domain.Recommendation{}, wrap("student", e)
	}
	if owner != in.CreatedBy {
		return domain.Recommendation{}, domain.ErrForbidden
	}
	if _, e = tx.ExecContext(ctx, "INSERT INTO recommendations(id,student_id,created_by,status,version,note,created_at,updated_at) VALUES(?,?,?,'draft',1,?,datetime('now'),datetime('now'))", in.ID, in.StudentID, in.CreatedBy, in.Note); e != nil {
		return domain.Recommendation{}, wrap("recommendation", e)
	}
	for _, it := range in.Items {
		var remain int
		if e = tx.QueryRowContext(ctx, "SELECT remaining FROM admission_plans WHERE id=?", it.AdmissionPlanID).Scan(&remain); e != nil {
			return domain.Recommendation{}, wrap("plan", e)
		}
		if remain <= 0 {
			return domain.Recommendation{}, fmt.Errorf("%w: quota exhausted", domain.ErrConflict)
		}
		if _, e = tx.ExecContext(ctx, "INSERT INTO recommendation_items(id,recommendation_id,admission_plan_id,career_path_id,priority) VALUES(lower(hex(randomblob(16))),?,?,?,?)", in.ID, it.AdmissionPlanID, it.CareerPathID, it.Priority); e != nil {
			return domain.Recommendation{}, wrap("item", e)
		}
	}
	if _, e = tx.ExecContext(ctx, "INSERT INTO audit_events(id,actor_id,entity_type,entity_id,action,result,request_id,detail,created_at) VALUES(lower(hex(randomblob(16))),?,'recommendation',?,'create','ok',?,?,datetime('now'))", in.CreatedBy, in.ID, requestID, in.Note); e != nil {
		return domain.Recommendation{}, wrap("audit", e)
	}
	if e = tx.Commit(); e != nil {
		return domain.Recommendation{}, wrap("commit", e)
	}
	return domain.Recommendation{ID: in.ID, StudentID: in.StudentID, CreatedBy: in.CreatedBy, Status: domain.RecommendationDraft, Version: 1, Note: in.Note}, nil
}
func (r RecommendationRepo) ByID(ctx context.Context, id string) (domain.Recommendation, error) {
	var x domain.Recommendation
	e := r.DB.QueryRowContext(ctx, "SELECT id,student_id,created_by,status,version,note FROM recommendations WHERE id=?", id).Scan(&x.ID, &x.StudentID, &x.CreatedBy, &x.Status, &x.Version, &x.Note)
	return x, wrap("recommendation", e)
}
func (r RecommendationRepo) Move(ctx context.Context, id string, from, to domain.RecommendationStatus, version int, actor, requestID string) error {
	// State transitions must not outlive the caller's request context.
	tx, e := r.DB.BeginTx(context.Background(), nil)
	if e != nil {
		return e
	}
	defer tx.Rollback()
	res, e := tx.ExecContext(context.Background(), "UPDATE recommendations SET status=?,version=version+1,updated_at=datetime('now'),submitted_at=CASE WHEN ?='submitted' THEN datetime('now') ELSE submitted_at END,approved_at=CASE WHEN ?='approved' THEN datetime('now') ELSE approved_at END WHERE id=? AND status=? AND version=?", to, to, to, id, from, version)
	if e != nil {
		return e
	}
	n, e := res.RowsAffected()
	if e != nil {
		return e
	}
	if n != 1 {
		return domain.ErrConflict
	}
	if _, e = tx.ExecContext(context.Background(), "INSERT INTO audit_events(id,actor_id,entity_type,entity_id,action,result,request_id,detail,created_at) VALUES(lower(hex(randomblob(16))),?,'recommendation',?,?,'ok',?, ?,datetime('now'))", actor, id, string(to), requestID, "state transition"); e != nil {
		return e
	}
	return tx.Commit()
}
func (r RecommendationRepo) List(ctx context.Context, student, status string) ([]domain.Recommendation, error) {
	q := "SELECT id,student_id,created_by,status,version,note FROM recommendations WHERE student_id=?"
	args := []any{student}
	if status != "" {
		q += " AND status=?"
		args = append(args, status)
	}
	q += " ORDER BY updated_at DESC"
	rows, e := r.DB.QueryContext(ctx, q, args...)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []domain.Recommendation{}
	for rows.Next() {
		var x domain.Recommendation
		if e = rows.Scan(&x.ID, &x.StudentID, &x.CreatedBy, &x.Status, &x.Version, &x.Note); e != nil {
			return nil, e
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

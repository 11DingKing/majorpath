package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/11DingKing/majorpath/internal/domain"
)

type StudentRepo struct{ DB *sql.DB }

func (r StudentRepo) Create(ctx context.Context, s domain.Student) error {
	_, e := r.DB.ExecContext(ctx, "INSERT INTO students(id,owner_id,name,graduation_year,region,version,created_at,updated_at) VALUES(?,?,?,?,?,1,datetime('now'),datetime('now'))", s.ID, s.OwnerID, s.Name, s.GraduationYear, s.Region)
	return wrap("create student", e)
}
func (r StudentRepo) ByID(ctx context.Context, id string) (domain.Student, error) {
	var s domain.Student
	e := r.DB.QueryRowContext(ctx, "SELECT id,owner_id,name,graduation_year,region,version FROM students WHERE id=?", id).Scan(&s.ID, &s.OwnerID, &s.Name, &s.GraduationYear, &s.Region, &s.Version)
	return s, wrap("find student", e)
}
func (r StudentRepo) AddInterest(ctx context.Context, student, tag string, weight float64) error {
	if tag == "" || weight <= 0 {
		return fmt.Errorf("%w: interest", domain.ErrInvalid)
	}
	_, e := r.DB.ExecContext(ctx, "INSERT INTO interests(id,student_id,tag,weight) VALUES(lower(hex(randomblob(16))),?,?,?)", student, tag, weight)
	return wrap("add interest", e)
}
func (r StudentRepo) List(ctx context.Context, owner string, limit, offset int) ([]domain.Student, error) {
	rows, e := r.DB.QueryContext(ctx, "SELECT id,owner_id,name,graduation_year,region,version FROM students WHERE owner_id=? ORDER BY created_at DESC LIMIT ? OFFSET ?", owner, limit, offset)
	if e != nil {
		return nil, wrap("list students", e)
	}
	defer rows.Close()
	out := make([]domain.Student, 0)
	for rows.Next() {
		var s domain.Student
		if e = rows.Scan(&s.ID, &s.OwnerID, &s.Name, &s.GraduationYear, &s.Region, &s.Version); e != nil {
			return nil, e
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

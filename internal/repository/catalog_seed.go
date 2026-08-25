package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type SeedRepo struct{ DB *sql.DB }

func (r SeedRepo) AddUniversity(ctx context.Context, id, name, province string) error {
	if id == "" || name == "" {
		return fmt.Errorf("university fields required")
	}
	_, e := r.DB.ExecContext(ctx, "INSERT INTO universities(id,name,province,created_at) VALUES(?,?,?,datetime('now'))", id, name, province)
	return wrap("university", e)
}
func (r SeedRepo) AddPlan(ctx context.Context, id, university, major string, year, quota int) error {
	if quota <= 0 {
		return fmt.Errorf("quota must be positive")
	}
	_, e := r.DB.ExecContext(ctx, "INSERT INTO admission_plans(id,university_id,major_id,year,quota,remaining) VALUES(?,?,?,?,?,?)", id, university, major, year, quota, quota)
	return wrap("plan", e)
}
func (r SeedRepo) AddPath(ctx context.Context, id, major, title, outlook string) error {
	if title == "" || outlook == "" {
		return fmt.Errorf("career path fields required")
	}
	_, e := r.DB.ExecContext(ctx, "INSERT INTO career_paths(id,major_id,title,outlook) VALUES(?,?,?,?)", id, major, title, outlook)
	return wrap("career path", e)
}
func (r SeedRepo) Counts(ctx context.Context) (map[string]int, error) {
	out := map[string]int{}
	for _, table := range []string{"users", "students", "majors", "universities", "admission_plans", "career_paths", "recommendations", "followups", "audit_events"} {
		var n int
		if e := r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+table).Scan(&n); e != nil {
			return nil, e
		}
		out[table] = n
	}
	return out, nil
}

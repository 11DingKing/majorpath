package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/majorpath/internal/domain"
)

type CatalogRepo struct{ DB *sql.DB }

func (r CatalogRepo) Seed(ctx context.Context) error {
	var n int
	if e := r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM majors").Scan(&n); e != nil {
		return e
	}
	if n > 0 {
		return nil
	}
	m := []struct{ c, n, f string }{{"CS101", "计算机科学与技术", "信息技术"}, {"DS202", "数据科学", "数据与人工智能"}, {"BIO303", "生物工程", "生命科学"}, {"ECO404", "经济学", "经济管理"}}
	for _, x := range m {
		if _, e := r.DB.ExecContext(ctx, "INSERT INTO majors(id,code,name,field,created_at) VALUES(lower(hex(randomblob(16))),?,?,?,datetime('now'))", x.c, x.n, x.f); e != nil {
			return e
		}
	}
	return nil
}
func (r CatalogRepo) Majors(ctx context.Context, field string) ([]domain.Major, error) {
	q := "SELECT id,code,name,field FROM majors"
	args := []any{}
	if field != "" {
		q += " WHERE field=?"
		args = append(args, field)
	}
	q += " ORDER BY name"
	rows, e := r.DB.QueryContext(ctx, q, args...)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []domain.Major{}
	for rows.Next() {
		var m domain.Major
		if e = rows.Scan(&m.ID, &m.Code, &m.Name, &m.Field); e != nil {
			return nil, e
		}
		out = append(out, m)
	}
	return out, rows.Err()
}
func (r CatalogRepo) Plans(ctx context.Context, major string, year int) ([]domain.AdmissionPlan, error) {
	rows, e := r.DB.QueryContext(ctx, "SELECT id,university_id,major_id,year,quota,remaining FROM admission_plans WHERE major_id=? AND year=? AND remaining>0 ORDER BY remaining DESC", major, year)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []domain.AdmissionPlan{}
	for rows.Next() {
		var p domain.AdmissionPlan
		if e = rows.Scan(&p.ID, &p.UniversityID, &p.MajorID, &p.Year, &p.Quota, &p.Remaining); e != nil {
			return nil, e
		}
		out = append(out, p)
	}
	return out, rows.Err()
}
func (r CatalogRepo) Path(ctx context.Context, major string) ([]domain.CareerPath, error) {
	rows, e := r.DB.QueryContext(ctx, "SELECT id,major_id,title,outlook FROM career_paths WHERE major_id=? ORDER BY title", major)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []domain.CareerPath{}
	for rows.Next() {
		var p domain.CareerPath
		if e = rows.Scan(&p.ID, &p.MajorID, &p.Title, &p.Outlook); e != nil {
			return nil, e
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

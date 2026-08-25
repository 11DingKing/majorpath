package repository

import (
	"context"
	"github.com/11DingKing/majorpath/internal/storage"
	"testing"
)

func TestCatalogSeedAndQueries(t *testing.T) {
	db, e := storage.Open(context.Background(), ":memory:")
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	r := CatalogRepo{DB: db}
	if e = r.Seed(context.Background()); e != nil {
		t.Fatal(e)
	}
	maj, e := r.Majors(context.Background(), "")
	if e != nil || len(maj) != 4 {
		t.Fatalf("majors=%d e=%v", len(maj), e)
	}
	filtered, e := r.Majors(context.Background(), "信息技术")
	if e != nil || len(filtered) != 1 {
		t.Fatal(filtered, e)
	}
	if e = r.Seed(context.Background()); e != nil {
		t.Fatal(e)
	}
	maj2, _ := r.Majors(context.Background(), "")
	if len(maj2) != 4 {
		t.Fatal("duplicate seed")
	}
}
func TestSeedRepoEntities(t *testing.T) {
	db, e := storage.Open(context.Background(), ":memory:")
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	r := SeedRepo{DB: db}
	if e = r.AddUniversity(context.Background(), "v", "大学", "浙江"); e != nil {
		t.Fatal(e)
	}
	if e = r.AddUniversity(context.Background(), "v", "重复", "浙江"); e == nil {
		t.Fatal("duplicate university")
	}
	db.Exec("INSERT INTO majors(id,code,name,field,created_at) VALUES('m','M','专业','领域','now')")
	if e = r.AddPlan(context.Background(), "p", "v", "m", 2026, 20); e != nil {
		t.Fatal(e)
	}
	if e = r.AddPath(context.Background(), "c", "m", "工程师", "需求增长"); e != nil {
		t.Fatal(e)
	}
	counts, e := r.Counts(context.Background())
	if e != nil || counts["universities"] != 1 || counts["admission_plans"] != 1 || counts["career_paths"] != 1 {
		t.Fatal(counts, e)
	}
}
func TestSeedRepoRejectsInvalid(t *testing.T) {
	db, e := storage.Open(context.Background(), ":memory:")
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	r := SeedRepo{DB: db}
	if e = r.AddUniversity(context.Background(), "", "大学", "P"); e == nil {
		t.Fatal("empty university")
	}
	if e = r.AddPlan(context.Background(), "p", "v", "m", 2026, 0); e == nil {
		t.Fatal("zero quota")
	}
	if e = r.AddPath(context.Background(), "c", "m", "", "outlook"); e == nil {
		t.Fatal("empty title")
	}
}
func TestReportSummariesEmpty(t *testing.T) {
	db, e := storage.Open(context.Background(), ":memory:")
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	r := ReportRepo{DB: db}
	v, e := r.Summaries(context.Background(), "none", 20, 0)
	if e != nil || len(v) != 0 {
		t.Fatal(v, e)
	}
}
func TestReportReserveUnknown(t *testing.T) {
	db, e := storage.Open(context.Background(), ":memory:")
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	r := ReportRepo{DB: db}
	if e = r.ReservePlan(context.Background(), "missing"); e == nil {
		t.Fatal("unknown plan reserved")
	}
}

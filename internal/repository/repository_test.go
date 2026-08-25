package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/11DingKing/majorpath/internal/domain"
	"github.com/11DingKing/majorpath/internal/storage"
	"testing"
	"time"
)

func repoDB(t *testing.T) *sql.DB {
	t.Helper()
	db, e := storage.Open(context.Background(), ":memory:")
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { db.Close() })
	return db
}
func seedUser(t *testing.T, db *sql.DB, id string) {
	t.Helper()
	_, e := db.Exec("INSERT INTO users(id,email,name,role,password_hash,created_at) VALUES(?,?,?,?,?,?)", id, id+"@e", id, "counselor", "p", "now")
	if e != nil {
		t.Fatal(e)
	}
}
func TestUserRepository(t *testing.T) {
	db := repoDB(t)
	seedUser(t, db, "u")
	r := UserRepo{DB: db}
	u, e := r.ByID(context.Background(), "u")
	if e != nil || u.Email != "u@e" {
		t.Fatal(u, e)
	}
	if _, e = r.ByID(context.Background(), "missing"); !errors.Is(e, domain.ErrNotFound) {
		t.Fatalf("e=%v", e)
	}
}
func TestStudentRepository(t *testing.T) {
	db := repoDB(t)
	seedUser(t, db, "u")
	r := StudentRepo{DB: db}
	s := domain.Student{ID: "s", OwnerID: "u", Name: "学生", GraduationYear: 2027, Region: "华东"}
	if e := r.Create(context.Background(), s); e != nil {
		t.Fatal(e)
	}
	got, e := r.ByID(context.Background(), "s")
	if e != nil || got.Name != "学生" {
		t.Fatal(got, e)
	}
	if e = r.AddInterest(context.Background(), "s", "数学", .8); e != nil {
		t.Fatal(e)
	}
	if e = r.AddInterest(context.Background(), "s", "数学", .7); e == nil {
		t.Fatal("duplicate interest")
	}
	list, e := r.List(context.Background(), "u", 10, 0)
	if e != nil || len(list) != 1 {
		t.Fatal(list, e)
	}
}
func TestStudentRepositoryValidation(t *testing.T) {
	db := repoDB(t)
	seedUser(t, db, "u")
	r := StudentRepo{DB: db}
	if e := r.AddInterest(context.Background(), "missing", "x", .4); e == nil {
		t.Fatal("missing student accepted")
	}
	if e := r.AddInterest(context.Background(), "missing", "", .4); e == nil {
		t.Fatal("empty tag accepted")
	}
}
func TestIdempotencyRepository(t *testing.T) {
	db := repoDB(t)
	seedUser(t, db, "u")
	r := IdempotencyRepo{DB: db}
	if e := r.Put(context.Background(), "k", "u", "h", "body"); e != nil {
		t.Fatal(e)
	}
	v, e := r.Get(context.Background(), "k", "u")
	if e != nil || v != "body" {
		t.Fatal(v, e)
	}
	if _, e = r.Get(context.Background(), "k", "other"); !errors.Is(e, domain.ErrNotFound) {
		t.Fatalf("e=%v", e)
	}
}
func TestFollowupRepositoryLifecycle(t *testing.T) {
	db := repoDB(t)
	seedUser(t, db, "u")
	db.Exec("INSERT INTO students(id,owner_id,name,graduation_year,region,version,created_at,updated_at) VALUES('s','u','S',2027,'r',1,'now','now')")
	db.Exec("INSERT INTO recommendations(id,student_id,created_by,status,version,note,created_at,updated_at) VALUES('r','s','u','approved',1,'n','now','now')")
	r := FollowupRepo{DB: db}
	due := time.Now().Add(-time.Minute)
	if e := r.Create(context.Background(), "f", "r", due); e != nil {
		t.Fatal(e)
	}
	f, e := r.Claim(context.Background(), time.Now(), time.Now().Add(time.Minute))
	if e != nil || f.Status != domain.FollowupRunning {
		t.Fatal(f, e)
	}
	if e = r.Complete(context.Background(), f.ID, true, ""); e != nil {
		t.Fatal(e)
	}
}
func TestPlanReserveRelease(t *testing.T) {
	db := repoDB(t)
	db.Exec("INSERT INTO majors(id,code,name,field,created_at) VALUES('m','M','Major','Field','now')")
	db.Exec("INSERT INTO universities(id,name,province,created_at) VALUES('v','Uni','P','now')")
	db.Exec("INSERT INTO admission_plans(id,university_id,major_id,year,quota,remaining) VALUES('p','v','m',2026,1,1)")
	r := ReportRepo{DB: db}
	if e := r.ReservePlan(context.Background(), "p"); e != nil {
		t.Fatal(e)
	}
	if e := r.ReservePlan(context.Background(), "p"); e == nil {
		t.Fatal("oversold")
	}
	if e := r.ReleasePlan(context.Background(), "p"); e != nil {
		t.Fatal(e)
	}
}
func TestRecommendationStateConflict(t *testing.T) {
	db := repoDB(t)
	seedUser(t, db, "u")
	db.Exec("INSERT INTO students(id,owner_id,name,graduation_year,region,version,created_at,updated_at) VALUES('s','u','S',2027,'r',1,'now','now')")
	r := RecommendationRepo{DB: db}
	x, e := r.Create(context.Background(), RecommendationInput{ID: "r", StudentID: "s", CreatedBy: "u", Note: "n", Items: []domain.RecommendationItem{{AdmissionPlanID: "p", CareerPathID: "c", Priority: 1}}}, "req")
	if e == nil {
		t.Fatal("missing plan accepted")
	}
	_ = x
}

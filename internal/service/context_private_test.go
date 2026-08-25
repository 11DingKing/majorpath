package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/11DingKing/majorpath/internal/domain"
	"github.com/11DingKing/majorpath/internal/repository"
	"github.com/11DingKing/majorpath/internal/storage"
)

func cancelledApprovalFixture(t *testing.T) (*sql.DB, domain.User, RecommendationService) {
	t.Helper()
	db, err := storage.Open(context.Background(), ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if _, err = db.Exec("INSERT INTO users(id,email,name,role,password_hash,created_at) VALUES('reviewer','reviewer@example.test','Reviewer','reviewer','hash','now')"); err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec("INSERT INTO students(id,owner_id,name,graduation_year,region,version,created_at,updated_at) VALUES('student','reviewer','Student',2027,'华东',1,'now','now')"); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"rec-cancel", "rec-live"} {
		if _, err = db.Exec("INSERT INTO recommendations(id,student_id,created_by,status,version,note,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?)", id, "student", "reviewer", "submitted", 1, id, "now", "now"); err != nil {
			t.Fatal(err)
		}
	}
	repo := repository.RecommendationRepo{DB: db}
	return db, domain.User{ID: "reviewer", Role: domain.RoleReviewer}, RecommendationService{Repo: repo}
}

func TestCancelledApprovalDoesNotCommit(t *testing.T) {
	db, reviewer, service := cancelledApprovalFixture(t)
	var before int
	if err := db.QueryRow("SELECT COUNT(*) FROM audit_events WHERE entity_id='rec-cancel'").Scan(&before); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := service.Approve(ctx, reviewer, "rec-cancel", 1); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled approval error = %v, want context.Canceled", err)
	}
	var status string
	if err := db.QueryRow("SELECT status FROM recommendations WHERE id='rec-cancel'").Scan(&status); err != nil {
		t.Fatal(err)
	}
	if status != string(domain.RecommendationSubmitted) {
		t.Fatalf("cancelled approval changed status to %q", status)
	}
	var after int
	if err := db.QueryRow("SELECT COUNT(*) FROM audit_events WHERE entity_id='rec-cancel'").Scan(&after); err != nil {
		t.Fatal(err)
	}
	if after != before {
		t.Fatalf("cancelled approval wrote %d audit events, want %d", after, before)
	}

	if err := service.Approve(context.Background(), reviewer, "rec-live", 1); err != nil {
		t.Fatalf("live approval failed: %v", err)
	}
	if err := db.QueryRow("SELECT status FROM recommendations WHERE id='rec-live'").Scan(&status); err != nil {
		t.Fatal(err)
	}
	if status != string(domain.RecommendationApproved) {
		t.Fatalf("live approval status = %q, want approved", status)
	}
	var liveAudit int
	if err := db.QueryRow("SELECT COUNT(*) FROM audit_events WHERE entity_id='rec-live' AND action='approved'").Scan(&liveAudit); err != nil {
		t.Fatal(err)
	}
	if liveAudit != 1 {
		t.Fatalf("live approval audit count = %d, want 1", liveAudit)
	}
}

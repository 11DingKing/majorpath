package service

import (
	"context"
	"errors"
	"testing"

	"github.com/11DingKing/majorpath/internal/domain"
	"github.com/11DingKing/majorpath/internal/repository"
	"github.com/11DingKing/majorpath/internal/storage"
)

func TestRecommendationQuotaFailureRollsBackWholeRequest(t *testing.T) {
	db, err := storage.Open(context.Background(), ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	for _, stmt := range []string{
		"INSERT INTO users(id,email,name,role,password_hash,created_at) VALUES('u-q','u-q@example.com','咨询师','counselor','p','now')",
		"INSERT INTO students(id,owner_id,name,graduation_year,region,version,created_at,updated_at) VALUES('s-q','u-q','学生',2027,'华东',1,'now','now')",
		"INSERT INTO majors(id,code,name,field,created_at) VALUES('m-q','MQ','专业','理工','now')",
		"INSERT INTO universities(id,name,province,created_at) VALUES('v-q','大学','省','now')",
		"INSERT INTO admission_plans(id,university_id,major_id,year,quota,remaining) VALUES('p-q','v-q','m-q',2026,1,1)",
		"INSERT INTO career_paths(id,major_id,title,outlook) VALUES('c-q','m-q','路径一','稳定')",
	} {
		if _, err := db.Exec(stmt); err != nil {
			t.Fatal(err)
		}
	}

	svc := RecommendationService{
		Repo:  repository.RecommendationRepo{DB: db},
		Audit: repository.AuditRepo{DB: db},
		Plans: repository.ReportRepo{DB: db},
	}
	_, err = svc.Create(context.Background(), domain.User{ID: "u-q", Role: domain.RoleCounselor}, "s-q", "两条方案", "req-q", []domain.RecommendationItem{
		{AdmissionPlanID: "p-q", CareerPathID: "c-q", Priority: 1},
		{AdmissionPlanID: "p-q", CareerPathID: "c-q", Priority: 2},
	})
	if !errors.Is(err, domain.ErrConflict) {
		t.Fatalf("expected quota conflict, got %v", err)
	}
	var recommendations, items, audits, remaining int
	if err := db.QueryRow("SELECT COUNT(*) FROM recommendations").Scan(&recommendations); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow("SELECT COUNT(*) FROM recommendation_items").Scan(&items); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow("SELECT COUNT(*) FROM audit_events WHERE entity_type='recommendation'").Scan(&audits); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow("SELECT remaining FROM admission_plans WHERE id='p-q'").Scan(&remaining); err != nil {
		t.Fatal(err)
	}
	if recommendations != 0 || items != 0 || audits != 0 || remaining != 1 {
		t.Fatalf("failed request leaked state: recommendations=%d items=%d audits=%d remaining=%d", recommendations, items, audits, remaining)
	}

	for _, stmt := range []string{
		"INSERT INTO admission_plans(id,university_id,major_id,year,quota,remaining) VALUES('p-q2','v-q','m-q',2027,1,1)",
		"INSERT INTO admission_plans(id,university_id,major_id,year,quota,remaining) VALUES('p-q3','v-q','m-q',2028,1,1)",
		"INSERT INTO career_paths(id,major_id,title,outlook) VALUES('c-q2','m-q','路径二','稳定')",
		"INSERT INTO career_paths(id,major_id,title,outlook) VALUES('c-q3','m-q','路径三','稳定')",
	} {
		if _, err := db.Exec(stmt); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := svc.Create(context.Background(), domain.User{ID: "u-q", Role: domain.RoleCounselor}, "s-q", "可用方案", "req-q2", []domain.RecommendationItem{
		{AdmissionPlanID: "p-q2", CareerPathID: "c-q2", Priority: 1},
		{AdmissionPlanID: "p-q3", CareerPathID: "c-q3", Priority: 2},
	}); err != nil {
		t.Fatalf("available plans should succeed: %v", err)
	}
	if err := db.QueryRow("SELECT COUNT(*) FROM recommendations").Scan(&recommendations); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow("SELECT COUNT(*) FROM recommendation_items").Scan(&items); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow("SELECT COUNT(*) FROM audit_events WHERE entity_type='recommendation'").Scan(&audits); err != nil {
		t.Fatal(err)
	}
	var left2, left3 int
	if err := db.QueryRow("SELECT remaining FROM admission_plans WHERE id='p-q2'").Scan(&left2); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow("SELECT remaining FROM admission_plans WHERE id='p-q3'").Scan(&left3); err != nil {
		t.Fatal(err)
	}
	if recommendations != 1 || items != 2 || audits != 1 || left2 != 0 || left3 != 0 {
		t.Fatalf("successful request mismatch: recommendations=%d items=%d audits=%d left2=%d left3=%d", recommendations, items, audits, left2, left3)
	}
}

package repository

import (
	"context"
	"github.com/11DingKing/majorpath/internal/storage"
	"testing"
	"time"
)

func TestCleanupExpired(t *testing.T) {
	db, e := storage.Open(context.Background(), ":memory:")
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	db.Exec("INSERT INTO users(id,email,name,role,password_hash,created_at) VALUES('u','u@e','U','counselor','p','now')")
	db.Exec("INSERT INTO sessions(id,user_id,token_hash,expires_at,created_at) VALUES('s','u','h','2000-01-01 00:00:00','now')")
	r := CleanupRepo{DB: db}
	n, e := r.PurgeExpiredSessions(context.Background(), time.Now())
	if e != nil || n != 1 {
		t.Fatal(n, e)
	}
}
func TestCleanupIdempotency(t *testing.T) {
	db, e := storage.Open(context.Background(), ":memory:")
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	db.Exec("INSERT INTO users(id,email,name,role,password_hash,created_at) VALUES('u','u@e','U','counselor','p','now')")
	db.Exec("INSERT INTO idempotency_keys(key,actor_id,request_hash,response_body,created_at,expires_at) VALUES('k','u','h','b','now','2000-01-01 00:00:00')")
	r := CleanupRepo{DB: db}
	n, e := r.PurgeIdempotency(context.Background(), time.Now())
	if e != nil || n != 1 {
		t.Fatal(n, e)
	}
}
func TestCleanupVacuum(t *testing.T) {
	db, e := storage.Open(context.Background(), ":memory:")
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	if e := (CleanupRepo{DB: db}).Vacuum(context.Background()); e != nil {
		t.Fatal(e)
	}
}

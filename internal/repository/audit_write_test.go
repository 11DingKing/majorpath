package repository

import (
	"context"
	"github.com/11DingKing/majorpath/internal/storage"
	"testing"
)

func TestAuditWriter(t *testing.T) {
	db, e := storage.Open(context.Background(), ":memory:")
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	db.Exec("INSERT INTO users(id,email,name,role,password_hash,created_at) VALUES('u','u@e','U','counselor','p','now')")
	w := AuditWriter{DB: db}
	if e = w.Write(context.Background(), "u", "student", "s", "create", "ok", "req", "created"); e != nil {
		t.Fatal(e)
	}
	n, e := w.Count(context.Background(), "student")
	if e != nil || n != 1 {
		t.Fatal(n, e)
	}
	if e = w.Write(context.Background(), "u", "", "s", "x", "ok", "req", "bad"); e == nil {
		t.Fatal("empty audit accepted")
	}
}
func TestAuditWriterContext(t *testing.T) {
	db, e := storage.Open(context.Background(), ":memory:")
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	ctx, c := context.WithCancel(context.Background())
	c()
	if e := (AuditWriter{DB: db}).Write(ctx, "u", "student", "s", "create", "ok", "req", "x"); e == nil {
		t.Fatal("cancel ignored")
	}
}

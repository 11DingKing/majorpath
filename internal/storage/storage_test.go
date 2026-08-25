package storage

import (
	"context"
	"database/sql"
	"github.com/11DingKing/majorpath/internal/repository"
	"os"
	"path/filepath"
	"testing"
)

func openTest(t *testing.T) (*sql.DB, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "data.db")
	db, e := Open(context.Background(), path)
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { db.Close(); os.Remove(path) })
	return db, path
}
func TestMigrationCreatesTables(t *testing.T) {
	db, _ := openTest(t)
	var n int
	if e := db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&n); e != nil {
		t.Fatal(e)
	}
	if n != 1 {
		t.Fatalf("migrations=%d", n)
	}
	rows, e := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'")
	if e != nil {
		t.Fatal(e)
	}
	defer rows.Close()
	seen := map[string]bool{}
	for rows.Next() {
		var name string
		rows.Scan(&name)
		seen[name] = true
	}
	for _, name := range []string{"users", "sessions", "students", "interests", "majors", "universities", "admission_plans", "career_paths", "recommendations", "recommendation_items", "followups", "audit_events", "idempotency_keys"} {
		if !seen[name] {
			t.Fatalf("missing %s", name)
		}
	}
}
func TestMigrationIdempotent(t *testing.T) {
	db, _ := openTest(t)
	if e := Migrate(context.Background(), db); e != nil {
		t.Fatal(e)
	}
	var n int
	db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&n)
	if n != 1 {
		t.Fatal(n)
	}
}
func TestForeignKeysEnabled(t *testing.T) {
	db, _ := openTest(t)
	if _, e := db.Exec("INSERT INTO students(id,owner_id,name,graduation_year,region,version,created_at,updated_at) VALUES('x','missing','n',2027,'r',1,'now','now')"); e == nil {
		t.Fatal("foreign key disabled")
	}
}
func TestRestartRecovery(t *testing.T) {
	db, path := openTest(t)
	if _, e := db.Exec("INSERT INTO users(id,email,name,role,password_hash,created_at) VALUES('u','u@e','U','counselor','p','now')"); e != nil {
		t.Fatal(e)
	}
	db.Close()
	db2, e := Open(context.Background(), path)
	if e != nil {
		t.Fatal(e)
	}
	defer db2.Close()
	var email string
	if e = db2.QueryRow("SELECT email FROM users WHERE id='u'").Scan(&email); e != nil || email != "u@e" {
		t.Fatalf("email=%s err=%v", email, e)
	}
}
func TestReady(t *testing.T) {
	db, _ := openTest(t)
	if e := Ready(context.Background(), db); e != nil {
		t.Fatal(e)
	}
}
func TestTxRunnerCommit(t *testing.T) {
	db, _ := openTest(t)
	r := repository.TxRunner{DB: db}
	if e := r.Run(context.Background(), func(tx *sql.Tx) error {
		_, e := tx.Exec("INSERT INTO users(id,email,name,role,password_hash,created_at) VALUES('u1','1@e','One','counselor','p','now')")
		return e
	}); e != nil {
		t.Fatal(e)
	}
	var n int
	db.QueryRow("SELECT COUNT(*) FROM users WHERE id='u1'").Scan(&n)
	if n != 1 {
		t.Fatal(n)
	}
}
func TestTxRunnerRollback(t *testing.T) {
	db, _ := openTest(t)
	r := repository.TxRunner{DB: db}
	if e := r.Run(context.Background(), func(tx *sql.Tx) error {
		tx.Exec("INSERT INTO users(id,email,name,role,password_hash,created_at) VALUES('u2','2@e','Two','counselor','p','now')")
		return context.Canceled
	}); e != context.Canceled {
		t.Fatal(e)
	}
	var n int
	db.QueryRow("SELECT COUNT(*) FROM users WHERE id='u2'").Scan(&n)
	if n != 0 {
		t.Fatal("rollback failed")
	}
}
func TestUniqueEmail(t *testing.T) {
	db, _ := openTest(t)
	q := "INSERT INTO users(id,email,name,role,password_hash,created_at) VALUES(?,?,?,?,?,?)"
	if _, e := db.Exec(q, "u1", "same@e", "One", "counselor", "p", "now"); e != nil {
		t.Fatal(e)
	}
	if _, e := db.Exec(q, "u2", "same@e", "Two", "counselor", "p", "now"); e == nil {
		t.Fatal("duplicate accepted")
	}
}
func TestIndexesExist(t *testing.T) {
	db, _ := openTest(t)
	rows, e := db.Query("SELECT name FROM sqlite_master WHERE type='index'")
	if e != nil {
		t.Fatal(e)
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		count++
	}
	if count < 4 {
		t.Fatalf("indexes=%d", count)
	}
}
func TestContextCancellation(t *testing.T) {
	db, _ := openTest(t)
	ctx, c := context.WithCancel(context.Background())
	c()
	if _, e := db.ExecContext(ctx, "SELECT 1"); e == nil {
		t.Fatal("cancel ignored")
	}
}
func TestMigrationVersionType(t *testing.T) {
	db, _ := openTest(t)
	var v int
	if e := db.QueryRow("SELECT version FROM schema_migrations").Scan(&v); e != nil || v != 1 {
		t.Fatalf("v=%d e=%v", v, e)
	}
}
func TestSchemaForeignKeyCascade(t *testing.T) {
	db, _ := openTest(t)
	db.Exec("INSERT INTO users(id,email,name,role,password_hash,created_at) VALUES('u3','3@e','Three','counselor','p','now')")
	db.Exec("INSERT INTO students(id,owner_id,name,graduation_year,region,version,created_at,updated_at) VALUES('s3','u3','S',2027,'r',1,'now','now')")
	db.Exec("INSERT INTO interests(id,student_id,tag,weight) VALUES('i3','s3','x',0.4)")
	db.Exec("DELETE FROM students WHERE id='s3'")
	var n int
	db.QueryRow("SELECT COUNT(*) FROM interests WHERE id='i3'").Scan(&n)
	if n != 0 {
		t.Fatal("cascade failed")
	}
}

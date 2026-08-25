package auth

import (
	"context"
	"database/sql"
	"errors"
	"github.com/11DingKing/majorpath/internal/clock"
	"github.com/11DingKing/majorpath/internal/domain"
	"github.com/11DingKing/majorpath/internal/repository"
	"github.com/11DingKing/majorpath/internal/storage"
	"testing"
	"time"
)

func authDB(t *testing.T) *sql.DB {
	db, e := storage.Open(context.Background(), ":memory:")
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { db.Close() })
	h := HashPassword("secret")
	db.Exec("INSERT INTO users(id,email,name,role,password_hash,created_at) VALUES('u','u@e','User','counselor',?,'now')", h)
	return db
}
func TestPasswordHash(t *testing.T) {
	h := HashPassword("secret")
	if h == "" || !CheckPassword(h, "secret") || CheckPassword(h, "other") {
		t.Fatal("password check")
	}
	if CheckPassword("bad", "secret") {
		t.Fatal("bad hash accepted")
	}
}
func TestTokenGeneration(t *testing.T) {
	a, e := NewToken()
	if e != nil || len(a) < 32 {
		t.Fatal(a, e)
	}
	b, _ := NewToken()
	if a == b {
		t.Fatal("tokens repeated")
	}
	if TokenPreview(a) == a || len(TokenPreview(a)) >= len(a) {
		t.Fatal("preview")
	}
	if TokenPreview("short") != "short" {
		t.Fatal("short preview")
	}
}
func TestLoginAndAuthenticate(t *testing.T) {
	db := authDB(t)
	s := Service{Users: repository.UserRepo{DB: db}, Sessions: repository.SessionRepo{DB: db}, Clock: clock.Fixed{T: time.Date(2035, 1, 1, 0, 0, 0, 0, time.UTC)}, TTL: time.Hour}
	token, u, e := s.Login(context.Background(), "u@e", "secret")
	if e != nil || u.ID != "u" || token == "" {
		t.Fatal(token, u, e)
	}
	got, e := s.Authenticate(context.Background(), token)
	if e != nil || got.ID != "u" {
		t.Fatal(got, e)
	}
	if e = s.Logout(context.Background(), token); e != nil {
		t.Fatal(e)
	}
	if _, e = s.Authenticate(context.Background(), token); !errors.Is(e, domain.ErrUnauthorized) {
		t.Fatalf("e=%v", e)
	}
}
func TestLoginRejectsPassword(t *testing.T) {
	db := authDB(t)
	s := Service{Users: repository.UserRepo{DB: db}, Sessions: repository.SessionRepo{DB: db}, Clock: clock.Real{}, TTL: time.Hour}
	if _, _, e := s.Login(context.Background(), "u@e", "wrong"); !errors.Is(e, domain.ErrUnauthorized) {
		t.Fatal(e)
	}
}
func TestLoginUnknownUser(t *testing.T) {
	db := authDB(t)
	s := Service{Users: repository.UserRepo{DB: db}, Sessions: repository.SessionRepo{DB: db}, Clock: clock.Real{}, TTL: time.Hour}
	if _, _, e := s.Login(context.Background(), "x@e", "secret"); !errors.Is(e, domain.ErrUnauthorized) {
		t.Fatal(e)
	}
}
func TestPolicy(t *testing.T) {
	p := Policy{}
	c := domain.User{ID: "c", Role: domain.RoleCounselor}
	r := domain.User{ID: "r", Role: domain.RoleReviewer}
	if e := p.CanCreateStudent(c); e != nil {
		t.Fatal(e)
	}
	if e := p.CanCreateStudent(r); e == nil {
		t.Fatal("reviewer create")
	}
	if e := p.CanSubmit(c, "c"); e != nil {
		t.Fatal(e)
	}
	if e := p.CanSubmit(c, "x"); e == nil {
		t.Fatal("foreign submit")
	}
	if e := p.CanReview(r); e != nil {
		t.Fatal(e)
	}
	if e := p.CanReview(c); e == nil {
		t.Fatal("counselor review")
	}
}
func TestSessionExpiry(t *testing.T) {
	db := authDB(t)
	now := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	s := Service{Users: repository.UserRepo{DB: db}, Sessions: repository.SessionRepo{DB: db}, Clock: clock.Fixed{T: now}, TTL: -time.Hour}
	token, _, e := s.Login(context.Background(), "u@e", "secret")
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.Authenticate(context.Background(), token); !errors.Is(e, domain.ErrUnauthorized) {
		t.Fatal(e)
	}
}

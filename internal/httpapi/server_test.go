package httpapi

import (
	"context"
	"database/sql"
	"github.com/11DingKing/majorpath/internal/auth"
	"github.com/11DingKing/majorpath/internal/clock"
	"github.com/11DingKing/majorpath/internal/repository"
	"github.com/11DingKing/majorpath/internal/service"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func testServer(t *testing.T) *Server {
	t.Helper()
	db, e := sql.Open("sqlite3", ":memory:")
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { db.Close() })
	return &Server{Auth: auth.Service{Users: repository.UserRepo{DB: db}, Sessions: repository.SessionRepo{DB: db}, Clock: clock.Fixed{T: time.Now()}, TTL: time.Hour}, Students: service.StudentService{Repo: repository.StudentRepo{DB: db}}, Catalog: service.CatalogService{Repo: repository.CatalogRepo{DB: db}}, Recommendations: service.RecommendationService{Repo: repository.RecommendationRepo{DB: db}}, Log: nil, Ready: func(context.Context) error { return nil }}
}
func TestHealth(t *testing.T) {
	s := testServer(t)
	r := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("status=%d", w.Code)
	}
}
func TestUnauthorized(t *testing.T) {
	s := testServer(t)
	r := httptest.NewRequest(http.MethodGet, "/v1/students", nil)
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, r)
	if w.Code != 401 {
		t.Fatalf("status=%d", w.Code)
	}
}
func TestLoginJSON(t *testing.T) {
	s := testServer(t)
	r := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader("{}"))
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, r)
	if w.Code != 401 {
		t.Fatalf("status=%d", w.Code)
	}
}

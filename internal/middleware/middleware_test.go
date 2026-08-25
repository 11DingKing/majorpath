package middleware

import (
	"github.com/11DingKing/majorpath/internal/auth"
	"github.com/11DingKing/majorpath/internal/domain"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestIDPreservesHeader(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(requestKey{}) != "req-fixed" {
			t.Fatal("missing id")
		}
		w.WriteHeader(204)
	})
	h := RequestID(next)
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("X-Request-ID", "req-fixed")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Header().Get("X-Request-ID") != "req-fixed" {
		t.Fatal(w.Header())
	}
}
func TestRequestIDGenerates(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(requestKey{}) == nil {
			t.Fatal("no id")
		}
		w.WriteHeader(204)
	})
	r := httptest.NewRequest(http.MethodGet, "/abc", nil)
	w := httptest.NewRecorder()
	RequestID(next).ServeHTTP(w, r)
	if w.Header().Get("X-Request-ID") == "" {
		t.Fatal("empty")
	}
}
func TestRecoverConvertsPanic(t *testing.T) {
	h := Recover(slog.Default(), http.HandlerFunc(func(http.ResponseWriter, *http.Request) { panic("boom") }))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
	if w.Code != 500 {
		t.Fatal(w.Code)
	}
}
func TestAuthRejectsMissing(t *testing.T) {
	a := auth.Service{}
	h := Auth(a, http.HandlerFunc(func(http.ResponseWriter, *http.Request) { t.Fatal("called") }))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
	if w.Code != 401 {
		t.Fatal(w.Code)
	}
}
func TestUserZero(t *testing.T) {
	u := User(httptest.NewRequest(http.MethodGet, "/", nil).Context())
	if u != (domain.User{}) {
		t.Fatal(u)
	}
}

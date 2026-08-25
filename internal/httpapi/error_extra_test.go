package httpapi

import (
	"errors"
	"github.com/11DingKing/majorpath/internal/domain"
	"net/http/httptest"
	"testing"
)

func TestErrorWrappingClassification(t *testing.T) {
	wrapped := errors.Join(errors.New("detail"), domain.ErrConflict)
	status, code := classify(wrapped)
	if status != 409 || code != "conflict" {
		t.Fatal(status, code)
	}
}
func TestHeadersNoCache(t *testing.T) {
	w := httptest.NewRecorder()
	setNoCache(w)
	if got := w.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatal(got)
	}
}
func TestMethodAllowed(t *testing.T) {
	w := httptest.NewRecorder()
	if methodAllowed(w, "POST") {
		t.Fatal("unexpected allowed")
	}
	if w.Header().Get("Allow") != "POST" {
		t.Fatal(w.Header())
	}
	w2 := httptest.NewRecorder()
	if !methodAllowed(w2, "") {
		t.Fatal("empty should allow")
	}
}

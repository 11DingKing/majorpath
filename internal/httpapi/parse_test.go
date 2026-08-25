package httpapi

import (
	"errors"
	"github.com/11DingKing/majorpath/internal/domain"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecodeBodyUnknownField(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"x","extra":1}`))
	var v struct {
		Name string `json:"name"`
	}
	if e := decodeBody(r, &v); e == nil {
		t.Fatal("unknown accepted")
	}
}
func TestDecodeBodyValid(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"x"}`))
	var v struct {
		Name string `json:"name"`
	}
	if e := decodeBody(r, &v); e != nil || v.Name != "x" {
		t.Fatal(v, e)
	}
}
func TestDecodeBodyNil(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	var v any
	if e := decodeBody(r, &v); !errors.Is(e, io.EOF) {
		t.Fatal(e)
	}
}
func TestBearer(t *testing.T) {
	cases := []struct {
		header string
		ok     bool
		token  string
	}{{"Bearer abc", true, "abc"}, {"bearer xyz", true, "xyz"}, {"", false, ""}, {"Token abc", false, ""}, {"Bearer", false, ""}}
	for _, c := range cases {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("Authorization", c.header)
		got, ok := bearer(r)
		if ok != c.ok || got != c.token {
			t.Fatalf("%q got %q %v", c.header, got, ok)
		}
	}
}
func TestContentType(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("Content-Type", "application/json; charset=utf-8")
	if !contentTypeJSON(r) {
		t.Fatal("json false")
	}
	r.Header.Set("Content-Type", "text/plain")
	if contentTypeJSON(r) {
		t.Fatal("plain true")
	}
}
func TestClassify(t *testing.T) {
	cases := []struct {
		e      error
		status int
		code   string
	}{{domain.ErrInvalid, 400, "invalid_request"}, {domain.ErrUnauthorized, 401, "unauthorized"}, {domain.ErrForbidden, 403, "forbidden"}, {domain.ErrNotFound, 404, "not_found"}, {domain.ErrConflict, 409, "conflict"}, {errors.New("x"), 500, "internal_error"}}
	for _, c := range cases {
		status, code := classify(c.e)
		if status != c.status || code != c.code {
			t.Fatalf("%v %d %s", c.e, status, code)
		}
	}
}
func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()
	writeError(w, 418, "teapot", "short")
	if w.Code != 418 || !strings.Contains(w.Body.String(), "teapot") {
		t.Fatal(w.Code, w.Body.String())
	}
}
func TestSetNoCache(t *testing.T) {
	w := httptest.NewRecorder()
	setNoCache(w)
	if w.Header().Get("Cache-Control") != "no-store" {
		t.Fatal(w.Header())
	}
}

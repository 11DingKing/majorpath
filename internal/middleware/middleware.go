package middleware

import (
	"context"
	"fmt"
	"github.com/11DingKing/majorpath/internal/auth"
	"github.com/11DingKing/majorpath/internal/domain"
	"log/slog"
	"net/http"
)

type key int

const userKey key = 1

func User(ctx context.Context) domain.User { u, _ := ctx.Value(userKey).(domain.User); return u }
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = fmt.Sprintf("req-%d", len(r.URL.Path))
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), requestKey{}, id)))
	})
}

type requestKey struct{}

func Auth(a auth.Service, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := r.Header.Get("Authorization")
		if len(t) < 8 {
			http.Error(w, "unauthorized", 401)
			return
		}
		u, e := a.Authenticate(r.Context(), t[7:])
		if e != nil {
			http.Error(w, "unauthorized", 401)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userKey, u)))
	})
}
func Recover(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if x := recover(); x != nil {
				log.Error("panic", "error", x)
				http.Error(w, "internal error", 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

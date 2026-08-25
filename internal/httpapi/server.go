package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/11DingKing/majorpath/internal/auth"
	"github.com/11DingKing/majorpath/internal/domain"
	"github.com/11DingKing/majorpath/internal/middleware"
	"github.com/11DingKing/majorpath/internal/query"
	"github.com/11DingKing/majorpath/internal/service"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	Auth            auth.Service
	Students        service.StudentService
	Catalog         service.CatalogService
	Recommendations service.RecommendationService
	Followups       service.FollowupService
	Log             *slog.Logger
	Ready           func(context.Context) error
}

func (s Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.health)
	mux.HandleFunc("GET /readyz", s.ready)
	mux.HandleFunc("POST /v1/auth/login", s.login)
	mux.HandleFunc("POST /v1/auth/logout", s.logout)
	mux.Handle("/v1/students", middleware.Auth(s.Auth, http.HandlerFunc(s.students)))
	mux.Handle("/v1/catalog/majors", middleware.Auth(s.Auth, http.HandlerFunc(s.majors)))
	mux.Handle("/v1/recommendations", middleware.Auth(s.Auth, http.HandlerFunc(s.recommendations)))
	return middleware.Recover(s.Log, middleware.RequestID(mux))
}
func write(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
func errCode(e error) int {
	switch {
	case errors.Is(e, domain.ErrUnauthorized), errors.Is(e, domain.ErrExpired):
		return 401
	case errors.Is(e, domain.ErrForbidden):
		return 403
	case errors.Is(e, domain.ErrNotFound):
		return 404
	case errors.Is(e, domain.ErrConflict):
		return 409
	case errors.Is(e, domain.ErrInvalid):
		return 400
	default:
		return 500
	}
}
func (s Server) health(w http.ResponseWriter, _ *http.Request) {
	write(w, 200, map[string]string{"status": "ok"})
}
func (s Server) ready(w http.ResponseWriter, r *http.Request) {
	if s.Ready != nil {
		if e := s.Ready(r.Context()); e != nil {
			write(w, 503, map[string]string{"status": "not_ready"})
			return
		}
	}
	write(w, 200, map[string]string{"status": "ready"})
}
func (s Server) login(w http.ResponseWriter, r *http.Request) {
	var in struct{ Email, Password string }
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		write(w, 400, map[string]string{"error": "invalid_json"})
		return
	}
	t, u, e := s.Auth.Login(r.Context(), in.Email, in.Password)
	if e != nil {
		write(w, errCode(e), map[string]string{"error": "unauthorized"})
		return
	}
	write(w, 200, map[string]any{"token": t, "user": u})
}
func (s Server) logout(w http.ResponseWriter, r *http.Request) {
	t := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if e := s.Auth.Logout(r.Context(), t); e != nil {
		write(w, errCode(e), map[string]string{"error": e.Error()})
		return
	}
	write(w, 204, nil)
}
func (s Server) students(w http.ResponseWriter, r *http.Request) {
	u := middleware.User(r.Context())
	switch r.Method {
	case http.MethodGet:
		p := query.Parse(r.URL.Query())
		v, e := s.Students.List(r.Context(), u, p.Limit, p.Offset)
		if e != nil {
			write(w, errCode(e), map[string]string{"error": e.Error()})
			return
		}
		write(w, 200, v)
	case http.MethodPost:
		var in struct {
			Name, Region   string
			GraduationYear int
		}
		if json.NewDecoder(r.Body).Decode(&in) != nil {
			write(w, 400, map[string]string{"error": "invalid_json"})
			return
		}
		v, e := s.Students.Create(r.Context(), u, in.Name, in.GraduationYear, in.Region)
		if e != nil {
			write(w, errCode(e), map[string]string{"error": e.Error()})
			return
		}
		write(w, 201, v)
	default:
		write(w, 405, nil)
	}
}
func (s Server) majors(w http.ResponseWriter, r *http.Request) {
	v, e := s.Catalog.Majors(r.Context(), r.URL.Query().Get("field"))
	if e != nil {
		write(w, 500, map[string]string{"error": e.Error()})
		return
	}
	write(w, 200, v)
}
func (s Server) recommendations(w http.ResponseWriter, r *http.Request) {
	u := middleware.User(r.Context())
	if r.Method == http.MethodGet {
		v, e := s.Recommendations.List(r.Context(), u, r.URL.Query().Get("student_id"), r.URL.Query().Get("status"))
		if e != nil {
			write(w, errCode(e), map[string]string{"error": e.Error()})
			return
		}
		write(w, 200, v)
		return
	}
	if r.Method != http.MethodPost {
		write(w, 405, nil)
		return
	}
	var in struct {
		StudentID, Note string
		Items           []domain.RecommendationItem
	}
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		write(w, 400, map[string]string{"error": "invalid_json"})
		return
	}
	v, e := s.Recommendations.Create(r.Context(), u, in.StudentID, in.Note, r.Header.Get("X-Request-ID"), in.Items)
	if e != nil {
		write(w, errCode(e), map[string]string{"error": e.Error()})
		return
	}
	write(w, 201, v)
}
func NewDefault(dbctx context.Context, a auth.Service, st service.StudentService, c service.CatalogService, r service.RecommendationService, f service.FollowupService, ready func(context.Context) error) Server {
	return Server{Auth: a, Students: st, Catalog: c, Recommendations: r, Followups: f, Log: slog.Default(), Ready: ready}
}

var _ = time.Second

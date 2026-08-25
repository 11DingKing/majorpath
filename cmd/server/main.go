package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"github.com/11DingKing/majorpath/internal/auth"
	"github.com/11DingKing/majorpath/internal/clock"
	"github.com/11DingKing/majorpath/internal/config"
	"github.com/11DingKing/majorpath/internal/httpapi"
	"github.com/11DingKing/majorpath/internal/repository"
	"github.com/11DingKing/majorpath/internal/service"
	"github.com/11DingKing/majorpath/internal/storage"
	"github.com/11DingKing/majorpath/internal/worker"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		os.Exit(0)
	}
	cfg := config.Load()
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	db, e := storage.Open(ctx, cfg.DatabasePath)
	if e != nil {
		panic(e)
	}
	defer db.Close()
	if e = seed(ctx, db); e != nil {
		panic(e)
	}
	clk := clock.Real{}
	users := repository.UserRepo{DB: db}
	sessions := repository.SessionRepo{DB: db}
	a := auth.Service{Users: users, Sessions: sessions, Clock: clk, TTL: cfg.SessionTTL}
	st := service.StudentService{Repo: repository.StudentRepo{DB: db}}
	cat := service.CatalogService{Repo: repository.CatalogRepo{DB: db}}
	rec := service.RecommendationService{Repo: repository.RecommendationRepo{DB: db}, Audit: repository.AuditRepo{DB: db}}
	fu := service.FollowupService{Repo: repository.FollowupRepo{DB: db}}
	srv := httpapi.NewDefault(ctx, a, st, cat, rec, fu, func(c context.Context) error { return db.PingContext(c) })
	w := worker.Worker{Followups: fu, Interval: cfg.WorkerInterval, Log: slog.Default()}
	go w.Run(ctx)
	h := &http.Server{Addr: ":" + cfg.Port, Handler: srv.Handler(), ReadHeaderTimeout: 5 * time.Second}
	go func() {
		<-ctx.Done()
		sh, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelShutdown()
		h.Shutdown(sh)
	}()
	slog.Info("majorpath listening", "addr", h.Addr)
	if e := h.ListenAndServe(); e != nil && e != http.ErrServerClosed {
		panic(e)
	}
}
func seed(ctx context.Context, db *sql.DB) error {
	h := sha256.Sum256([]byte("password"))
	p := hex.EncodeToString(h[:])
	if _, e := db.ExecContext(ctx, "INSERT OR IGNORE INTO users(id,email,name,role,password_hash,created_at) VALUES('u-counselor','counselor@example.com','咨询师','counselor',?,datetime('now'))", p); e != nil {
		return e
	}
	if _, e := db.ExecContext(ctx, "INSERT OR IGNORE INTO users(id,email,name,role,password_hash,created_at) VALUES('u-reviewer','reviewer@example.com','审核员','reviewer',?,datetime('now'))", p); e != nil {
		return e
	}
	return repository.CatalogRepo{DB: db}.Seed(ctx)
}

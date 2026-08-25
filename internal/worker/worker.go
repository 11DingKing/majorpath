package worker

import (
	"context"
	"errors"
	"github.com/11DingKing/majorpath/internal/domain"
	"github.com/11DingKing/majorpath/internal/service"
	"log/slog"
	"time"
)

type Worker struct {
	Followups service.FollowupService
	Interval  time.Duration
	Log       *slog.Logger
}

func (w Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			w.once(ctx, t)
		}
	}
}
func (w Worker) once(ctx context.Context, now time.Time) {
	f, e := w.Followups.Claim(ctx, now)
	if e != nil {
		if !errors.Is(e, domain.ErrNotFound) {
			w.Log.Error("followup claim", "error", e)
		}
		return
	}
	if e = w.Followups.Complete(ctx, f.ID, true, ""); e != nil {
		w.Log.Error("followup complete", "error", e)
	}
}

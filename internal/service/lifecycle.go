package service

import (
	"context"
	"fmt"
	"github.com/11DingKing/majorpath/internal/domain"
	"time"
)

type Lifecycle struct{ Clock func() time.Time }

func (l Lifecycle) Deadline(due time.Time) error {
	now := time.Now
	if l.Clock != nil {
		now = l.Clock
	}
	if now().After(due) {
		return fmt.Errorf("%w: deadline passed", domain.ErrExpired)
	}
	return nil
}
func (l Lifecycle) RequireActive(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
func (l Lifecycle) NextFollowupStatus(current domain.FollowupStatus, success bool) (domain.FollowupStatus, error) {
	switch current {
	case domain.FollowupPending:
		return domain.FollowupRunning, nil
	case domain.FollowupRunning:
		if success {
			return domain.FollowupDone, nil
		}
		return domain.FollowupFailed, nil
	case domain.FollowupFailed:
		return domain.FollowupRunning, nil
	default:
		return current, fmt.Errorf("%w: lifecycle", domain.ErrConflict)
	}
}
func (l Lifecycle) RetryAllowed(attempts, max int) bool { return attempts < max && attempts >= 0 }

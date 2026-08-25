package service

import (
	"context"
	"github.com/11DingKing/majorpath/internal/domain"
	"github.com/11DingKing/majorpath/internal/repository"
	"time"
)

type FollowupService struct{ Repo repository.FollowupRepo }

func (s FollowupService) Schedule(ctx context.Context, u domain.User, id string, due time.Time) error {
	if u.Role != domain.RoleCounselor {
		return domain.ErrForbidden
	}
	return s.Repo.Create(ctx, id, id, due)
}
func (s FollowupService) Claim(ctx context.Context, now time.Time) (domain.Followup, error) {
	return s.Repo.Claim(ctx, now, now.Add(2*time.Minute))
}
func (s FollowupService) Complete(ctx context.Context, id string, ok bool, msg string) error {
	return s.Repo.Complete(ctx, id, ok, msg)
}

package service

import (
	"context"
	"github.com/11DingKing/majorpath/internal/domain"
	"github.com/11DingKing/majorpath/internal/repository"
)

type ReportingService struct{ Repo repository.ReportRepo }

func (s ReportingService) Summaries(ctx context.Context, u domain.User, limit, offset int) ([]repository.RecommendationSummary, error) {
	if u.ID == "" {
		return nil, domain.ErrUnauthorized
	}
	return s.Repo.Summaries(ctx, u.ID, limit, offset)
}
func (s ReportingService) Reserve(ctx context.Context, u domain.User, plan string) error {
	if u.Role != domain.RoleCounselor {
		return domain.ErrForbidden
	}
	return s.Repo.ReservePlan(ctx, plan)
}
func (s ReportingService) Release(ctx context.Context, u domain.User, plan string) error {
	if u.Role != domain.RoleReviewer {
		return domain.ErrForbidden
	}
	return s.Repo.ReleasePlan(ctx, plan)
}

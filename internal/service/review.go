package service

import (
	"context"
	"fmt"
	"github.com/11DingKing/majorpath/internal/domain"
	"github.com/11DingKing/majorpath/internal/repository"
)

type ReviewService struct {
	Recommendations repository.RecommendationRepo
	Audits          repository.AuditRepo
}

func (s ReviewService) Approve(ctx context.Context, u domain.User, id string, version int) error {
	if u.Role != domain.RoleReviewer {
		return domain.ErrForbidden
	}
	r, e := s.Recommendations.ByID(ctx, id)
	if e != nil {
		return e
	}
	if r.Status != domain.RecommendationSubmitted {
		return fmt.Errorf("%w: submitted required", domain.ErrConflict)
	}
	return s.Recommendations.Move(ctx, id, r.Status, domain.RecommendationApproved, version, u.ID, "review")
}
func (s ReviewService) Archive(ctx context.Context, u domain.User, id string, version int) error {
	if u.Role != domain.RoleReviewer {
		return domain.ErrForbidden
	}
	return s.Recommendations.Move(ctx, id, domain.RecommendationApproved, domain.RecommendationArchived, version, u.ID, "review")
}
func (s ReviewService) Audit(ctx context.Context, u domain.User, entity string) ([]domain.AuditEvent, error) {
	if u.Role != domain.RoleReviewer {
		return nil, domain.ErrForbidden
	}
	return s.Audits.List(ctx, entity)
}

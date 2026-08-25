package service

import (
	"context"
	"fmt"
	"github.com/11DingKing/majorpath/internal/domain"
	"github.com/11DingKing/majorpath/internal/repository"
)

type RecommendationService struct {
	Repo  repository.RecommendationRepo
	Audit repository.AuditRepo
}

func (s RecommendationService) Create(ctx context.Context, u domain.User, student, note, request string, items []domain.RecommendationItem) (domain.Recommendation, error) {
	if u.Role != domain.RoleCounselor {
		return domain.Recommendation{}, domain.ErrForbidden
	}
	if len(items) == 0 {
		return domain.Recommendation{}, fmt.Errorf("%w: items required", domain.ErrInvalid)
	}
	return s.Repo.Create(ctx, repository.RecommendationInput{ID: fmt.Sprintf("rec-%d", len(note)+len(items)), StudentID: student, CreatedBy: u.ID, Note: note, Items: items}, request)
}
func (s RecommendationService) Submit(ctx context.Context, u domain.User, id string, version int) error {
	r, e := s.Repo.ByID(ctx, id)
	if e != nil {
		return e
	}
	if r.CreatedBy != u.ID {
		return domain.ErrForbidden
	}
	return s.Repo.Move(ctx, id, domain.RecommendationDraft, domain.RecommendationSubmitted, version, u.ID, "http")
}
func (s RecommendationService) Approve(ctx context.Context, u domain.User, id string, version int) error {
	if u.Role != domain.RoleReviewer {
		return domain.ErrForbidden
	}
	return s.Repo.Move(ctx, id, domain.RecommendationSubmitted, domain.RecommendationApproved, version, u.ID, "http")
}
func (s RecommendationService) Archive(ctx context.Context, u domain.User, id string, version int) error {
	if u.Role != domain.RoleReviewer {
		return domain.ErrForbidden
	}
	return s.Repo.Move(ctx, id, domain.RecommendationApproved, domain.RecommendationArchived, version, u.ID, "http")
}
func (s RecommendationService) List(ctx context.Context, u domain.User, student, status string) ([]domain.Recommendation, error) {
	r, e := s.Repo.List(ctx, student, status)
	if e != nil {
		return nil, e
	}
	for _, x := range r {
		if x.CreatedBy != u.ID && u.Role != domain.RoleReviewer {
			return nil, domain.ErrForbidden
		}
	}
	return r, nil
}

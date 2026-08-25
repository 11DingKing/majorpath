package service

import (
	"context"
	"fmt"
	"github.com/11DingKing/majorpath/internal/domain"
	"github.com/11DingKing/majorpath/internal/repository"
)

type ConsistencyService struct {
	Students repository.StudentRepo
	Plans    repository.ReportRepo
}

func (s ConsistencyService) CheckStudentOwner(ctx context.Context, u domain.User, id string) error {
	x, e := s.Students.ByID(ctx, id)
	if e != nil {
		return e
	}
	if x.OwnerID != u.ID {
		return fmt.Errorf("%w: owner mismatch", domain.ErrForbidden)
	}
	return nil
}
func (s ConsistencyService) ReserveForStudent(ctx context.Context, u domain.User, student, plan string) error {
	if e := s.CheckStudentOwner(ctx, u, student); e != nil {
		return e
	}
	return s.Plans.ReservePlan(ctx, plan)
}
func (s ConsistencyService) ReleaseForReviewer(ctx context.Context, u domain.User, plan string) error {
	if u.Role != domain.RoleReviewer {
		return domain.ErrForbidden
	}
	return s.Plans.ReleasePlan(ctx, plan)
}
func (s ConsistencyService) ValidateRecommendation(ctx context.Context, r domain.Recommendation) error {
	if r.ID == "" || r.StudentID == "" {
		return fmt.Errorf("%w: recommendation identity", domain.ErrInvalid)
	}
	if r.Version < 1 {
		return fmt.Errorf("%w: recommendation version", domain.ErrInvalid)
	}
	return EnsureContext(ctx)
}

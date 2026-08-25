package audit

import (
	"context"
	"github.com/11DingKing/majorpath/internal/domain"
	"github.com/11DingKing/majorpath/internal/repository"
)

type Service struct{ Repo repository.AuditRepo }

func (s Service) List(ctx context.Context, u domain.User, entity string) ([]domain.AuditEvent, error) {
	if u.Role != domain.RoleReviewer {
		return nil, domain.ErrForbidden
	}
	return s.Repo.List(ctx, entity)
}

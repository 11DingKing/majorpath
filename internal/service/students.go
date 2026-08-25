package service

import (
	"context"
	"fmt"
	"github.com/11DingKing/majorpath/internal/domain"
	"github.com/11DingKing/majorpath/internal/repository"
)

type StudentService struct{ Repo repository.StudentRepo }

func (s StudentService) Create(ctx context.Context, u domain.User, name string, year int, region string) (domain.Student, error) {
	if u.Role != domain.RoleCounselor {
		return domain.Student{}, domain.ErrForbidden
	}
	if name == "" || year < 2020 {
		return domain.Student{}, fmt.Errorf("%w: student fields", domain.ErrInvalid)
	}
	x := domain.Student{ID: fmt.Sprintf("stu-%d", len(name)+year), OwnerID: u.ID, Name: name, GraduationYear: year, Region: region, Version: 1}
	return x, s.Repo.Create(ctx, x)
}
func (s StudentService) AddInterest(ctx context.Context, u domain.User, id, tag string, w float64) error {
	x, e := s.Repo.ByID(ctx, id)
	if e != nil {
		return e
	}
	if x.OwnerID != u.ID {
		return domain.ErrForbidden
	}
	return s.Repo.AddInterest(ctx, id, tag, w)
}
func (s StudentService) List(ctx context.Context, u domain.User, limit, offset int) ([]domain.Student, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.Repo.List(ctx, u.ID, limit, offset)
}

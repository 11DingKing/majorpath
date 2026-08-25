package service

import (
	"context"
	"github.com/11DingKing/majorpath/internal/domain"
	"github.com/11DingKing/majorpath/internal/repository"
)

type CatalogService struct{ Repo repository.CatalogRepo }

func (s CatalogService) Majors(ctx context.Context, field string) ([]domain.Major, error) {
	return s.Repo.Majors(ctx, field)
}
func (s CatalogService) Plans(ctx context.Context, major string, year int) ([]domain.AdmissionPlan, error) {
	return s.Repo.Plans(ctx, major, year)
}
func (s CatalogService) Paths(ctx context.Context, major string) ([]domain.CareerPath, error) {
	return s.Repo.Path(ctx, major)
}

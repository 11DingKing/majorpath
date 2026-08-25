package service

import (
	"context"
	"github.com/11DingKing/majorpath/internal/domain"
	"sort"
)

type StudentView struct {
	Student       domain.Student
	InterestCount int
}

func SortStudents(v []domain.Student, descending bool) {
	sort.SliceStable(v, func(i, j int) bool {
		if descending {
			return v[i].GraduationYear > v[j].GraduationYear
		}
		return v[i].GraduationYear < v[j].GraduationYear
	})
}
func FilterRegion(v []domain.Student, region string) []domain.Student {
	out := make([]domain.Student, 0, len(v))
	for _, x := range v {
		if region == "" || x.Region == region {
			out = append(out, x)
		}
	}
	return out
}
func EnsureContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

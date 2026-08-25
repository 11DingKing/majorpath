package service

import (
	"context"
	"github.com/11DingKing/majorpath/internal/domain"
	"testing"
)

type fakeStudentRepo struct{ created []domain.Student }

func (f *fakeStudentRepo) Create(_ context.Context, s domain.Student) error {
	f.created = append(f.created, s)
	return nil
}
func (f *fakeStudentRepo) ByID(context.Context, string) (domain.Student, error) {
	return domain.Student{ID: "s", OwnerID: "u"}, nil
}
func (f *fakeStudentRepo) AddInterest(context.Context, string, string, float64) error { return nil }
func (f *fakeStudentRepo) List(context.Context, string, int, int) ([]domain.Student, error) {
	return f.created, nil
}
func TestSortStudents(t *testing.T) {
	v := []domain.Student{{Name: "a", GraduationYear: 2028}, {Name: "b", GraduationYear: 2026}, {Name: "c", GraduationYear: 2027}}
	SortStudents(v, false)
	if v[0].GraduationYear != 2026 || v[2].GraduationYear != 2028 {
		t.Fatal(v)
	}
	SortStudents(v, true)
	if v[0].GraduationYear != 2028 {
		t.Fatal(v)
	}
}
func TestFilterRegion(t *testing.T) {
	v := []domain.Student{{Region: "华东"}, {Region: "华南"}, {Region: "华东"}}
	if len(FilterRegion(v, "华东")) != 2 {
		t.Fatal("filter")
	}
	if len(FilterRegion(v, "")) != 3 {
		t.Fatal("all")
	}
}
func TestEnsureContext(t *testing.T) {
	ctx, c := context.WithCancel(context.Background())
	c()
	if EnsureContext(ctx) == nil {
		t.Fatal("cancel ignored")
	}
	if EnsureContext(context.Background()) != nil {
		t.Fatal("live context rejected")
	}
}
func TestStudentViewIsValue(t *testing.T) {
	v := StudentView{Student: domain.Student{Name: "test"}, InterestCount: 2}
	if v.Student.Name != "test" || v.InterestCount != 2 {
		t.Fatal(v)
	}
}

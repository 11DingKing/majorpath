package domain

import "testing"

func TestPageInfo(t *testing.T) {
	cases := []struct {
		l, o, total int
		more        bool
	}{{20, 0, 100, true}, {20, 80, 100, false}, {0, 0, 5, false}, {10, -1, 15, true}}
	for _, c := range cases {
		p := NewPageInfo(c.l, c.o, c.total)
		if p.HasMore != c.more {
			t.Fatalf("%+v", p)
		}
		if p.Limit < 1 || p.Offset < 0 {
			t.Fatalf("invalid page %+v", p)
		}
	}
}
func TestSortDirection(t *testing.T) {
	if NormalizeSort("desc") != SortDesc {
		t.Fatal("desc")
	}
	if NormalizeSort("asc") != SortAsc {
		t.Fatal("asc")
	}
	if NormalizeSort("random") != SortAsc {
		t.Fatal("fallback")
	}
}
func TestStatusStrings(t *testing.T) {
	if string(FollowupPending) != "pending" || string(FollowupDone) != "done" {
		t.Fatal("status")
	}
	if string(RecommendationSubmitted) != "submitted" {
		t.Fatal("recommendation")
	}
}
func TestErrorValuesDistinct(t *testing.T) {
	values := []error{ErrNotFound, ErrConflict, ErrForbidden, ErrInvalid, ErrExpired, ErrUnauthorized, ErrUnavailable}
	for i := range values {
		for j := range values {
			if i != j && values[i] == values[j] {
				t.Fatal("duplicate error")
			}
		}
	}
}
func TestStudentZeroValue(t *testing.T) {
	var s Student
	if s.ID != "" || s.Version != 0 {
		t.Fatal(s)
	}
}

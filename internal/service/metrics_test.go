package service

import (
	"context"
	"errors"
	"testing"
)

func TestMetricsValidation(t *testing.T) {
	good := Metrics{Students: 10, Recommendations: 4, Approved: 2, Followups: 3}
	if e := good.Validate(); e != nil {
		t.Fatal(e)
	}
	bad := []Metrics{{Students: -1}, {Recommendations: -1}, {Recommendations: 2, Approved: 3}, {Followups: -1}}
	for _, m := range bad {
		if e := m.Validate(); e == nil {
			t.Fatalf("accepted %+v", m)
		}
	}
}
func TestMetricsRate(t *testing.T) {
	if got := (Metrics{}).ApprovalRate(); got != 0 {
		t.Fatal(got)
	}
	if got := (Metrics{Recommendations: 4, Approved: 1}).ApprovalRate(); got != .25 {
		t.Fatal(got)
	}
}
func TestMetricsMerge(t *testing.T) {
	a := Metrics{Students: 1, Recommendations: 2, Approved: 1, Followups: 3}
	b := Metrics{Students: 4, Recommendations: 5, Approved: 2, Followups: 1}
	if got := a.Merge(b); got != (Metrics{5, 7, 3, 4}) {
		t.Fatal(got)
	}
}
func TestMetricsContext(t *testing.T) {
	m := Metrics{Students: 1}
	if e := m.CheckContext(context.Background()); e != nil {
		t.Fatal(e)
	}
	ctx, c := context.WithCancel(context.Background())
	c()
	if e := m.CheckContext(ctx); !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
}

func TestMetricsBoundaryValues(t *testing.T) {
	cases := []Metrics{
		{Students: 0, Recommendations: 0, Approved: 0, Followups: 0},
		{Students: 1, Recommendations: 1, Approved: 1, Followups: 1},
		{Students: 100, Recommendations: 50, Approved: 50, Followups: 20},
	}
	for _, item := range cases {
		if err := item.Validate(); err != nil {
			t.Fatalf("boundary %+v: %v", item, err)
		}
	}
}

func TestMetricsMergePreservesCounts(t *testing.T) {
	base := Metrics{Students: 2, Recommendations: 3, Approved: 1, Followups: 4}
	zero := Metrics{}
	if got := base.Merge(zero); got != base {
		t.Fatalf("zero merge changed metrics: %+v", got)
	}
}

func TestMetricsApprovalRateFullRange(t *testing.T) {
	for approved := 0; approved <= 10; approved++ {
		m := Metrics{Recommendations: 10, Approved: approved}
		got := m.ApprovalRate()
		if got < 0 || got > 1 {
			t.Fatalf("rate out of range for %d: %v", approved, got)
		}
	}
}

func TestMetricsStringSafe(t *testing.T) {
	m := Metrics{Students: 1, Recommendations: 1, Approved: 1, Followups: 1}
	if err := m.Validate(); err != nil {
		t.Fatal(err)
	}
}

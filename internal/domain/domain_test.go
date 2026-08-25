package domain

import (
	"testing"
	"time"
)

func TestRecommendationTransitions(t *testing.T) {
	cases := []struct {
		from, to RecommendationStatus
		ok       bool
	}{{RecommendationDraft, RecommendationSubmitted, true}, {RecommendationSubmitted, RecommendationApproved, true}, {RecommendationApproved, RecommendationArchived, true}, {RecommendationDraft, RecommendationApproved, false}, {RecommendationArchived, RecommendationDraft, false}}
	for _, c := range cases {
		if got := c.from.CanMove(c.to); got != c.ok {
			t.Fatalf("%s -> %s got %v", c.from, c.to, got)
		}
	}
}
func TestStudentProfileValidation(t *testing.T) {
	now := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	valid := StudentProfile{Name: "林晓", GraduationYear: 2027, Region: "华东", Interests: map[string]float64{"算法": .8}}
	if e := valid.Validate(now); e != nil {
		t.Fatal(e)
	}
	cases := []StudentProfile{{GraduationYear: 2027, Region: "华东"}, {Name: "a", GraduationYear: 2010, Region: "华东"}, {Name: "a", GraduationYear: 2027, Region: ""}, {Name: "a", GraduationYear: 2027, Region: "x", Interests: map[string]float64{"a": 0}}}
	for i, c := range cases {
		if e := c.Validate(now); e == nil {
			t.Fatalf("case %d accepted", i)
		}
	}
}
func TestValidateTransition(t *testing.T) {
	if e := ValidateTransition(RecommendationDraft, RecommendationSubmitted); e != nil {
		t.Fatal(e)
	}
	if e := ValidateTransition(RecommendationDraft, RecommendationDraft); e == nil {
		t.Fatal("same state accepted")
	}
	if e := ValidateTransition(RecommendationArchived, RecommendationDraft); e == nil {
		t.Fatal("backward transition accepted")
	}
}
func TestFollowupValidation(t *testing.T) {
	for _, s := range []FollowupStatus{FollowupPending, FollowupRunning, FollowupDone, FollowupFailed} {
		if e := ValidateFollowup(s); e != nil {
			t.Fatal(e)
		}
	}
	if e := ValidateFollowup("unknown"); e == nil {
		t.Fatal("unknown status accepted")
	}
}
func TestRoleValues(t *testing.T) {
	if RoleCounselor == RoleReviewer {
		t.Fatal("roles overlap")
	}
	if RecommendationDraft == RecommendationApproved {
		t.Fatal("states overlap")
	}
}

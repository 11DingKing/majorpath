package domain

import "testing"

func TestAdditionalCapabilities(t *testing.T) {
	c := User{ID: "c", Role: RoleCounselor}
	r := User{ID: "r", Role: RoleReviewer}
	if !c.CanAssignPlan() {
		t.Fatal("counselor cannot assign")
	}
	if r.CanAssignPlan() {
		t.Fatal("reviewer assigned")
	}
	if !c.CanManageSessions() || !r.CanManageSessions() {
		t.Fatal("session management")
	}
	if !r.CanInspectMetrics() || c.CanInspectMetrics() {
		t.Fatal("metrics")
	}
	if !c.IsKnown() || !r.IsKnown() || (User{Role: "unknown"}).IsKnown() {
		t.Fatal("known role")
	}
}

package domain

import "testing"

func TestUserRoleCapabilities(t *testing.T) {
	c := User{ID: "c", Role: RoleCounselor}
	r := User{ID: "r", Role: RoleReviewer}
	if !c.IsCounselor() || c.IsReviewer() {
		t.Fatal("counselor role")
	}
	if !r.IsReviewer() || r.IsCounselor() {
		t.Fatal("reviewer role")
	}
	if !c.CanReadStudent("c") || c.CanReadStudent("other") {
		t.Fatal("counselor read")
	}
	if !r.CanReadStudent("other") || !c.CanEditStudent("c") || c.CanEditStudent("other") {
		t.Fatal("ownership")
	}
	if !r.CanApprove() || c.CanApprove() {
		t.Fatal("approval")
	}
	if !r.CanViewAudit() || c.CanViewAudit() {
		t.Fatal("audit")
	}
	if !c.CanScheduleFollowup("c") || c.CanScheduleFollowup("r") {
		t.Fatal("followup")
	}
	if !r.CanArchive() || c.CanArchive() {
		t.Fatal("archive")
	}
	if c.RoleName() != "counselor" || !r.CanManageCatalog() || c.CanManageCatalog() {
		t.Fatal("catalog")
	}
	if !c.CanExport() || !r.CanExport() || !c.CanDeleteStudent("c") || c.CanDeleteStudent("r") {
		t.Fatal("export/delete")
	}
	if !c.Same(User{ID: "c"}) || c.Same(r) || (User{}).Same(c) {
		t.Fatal("identity")
	}
}

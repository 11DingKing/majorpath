package domain

func (u User) IsCounselor() bool                     { return u.Role == RoleCounselor }
func (u User) IsReviewer() bool                      { return u.Role == RoleReviewer }
func (u User) CanReadStudent(owner string) bool      { return u.IsReviewer() || u.ID == owner }
func (u User) CanEditStudent(owner string) bool      { return u.IsCounselor() && u.ID == owner }
func (u User) CanApprove() bool                      { return u.IsReviewer() }
func (u User) CanViewAudit() bool                    { return u.IsReviewer() }
func (u User) CanScheduleFollowup(owner string) bool { return u.IsCounselor() && u.ID == owner }
func (u User) CanArchive() bool                      { return u.IsReviewer() }
func (u User) RoleName() string                      { return string(u.Role) }
func (u User) CanManageCatalog() bool                { return u.IsReviewer() }
func (u User) CanExport() bool                       { return u.IsReviewer() || u.IsCounselor() }
func (u User) CanDeleteStudent(owner string) bool    { return u.IsCounselor() && u.ID == owner }
func (u User) Same(other User) bool                  { return u.ID != "" && u.ID == other.ID }
func (u User) CanAssignPlan() bool                   { return u.IsCounselor() }
func (u User) CanManageSessions() bool               { return u.ID != "" }
func (u User) CanInspectMetrics() bool               { return u.IsReviewer() }
func (u User) IsKnown() bool                         { return u.IsCounselor() || u.IsReviewer() }

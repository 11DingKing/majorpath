package domain

import "errors"

var (
	ErrNotFound     = errors.New("resource not found")
	ErrConflict     = errors.New("resource conflict")
	ErrForbidden    = errors.New("forbidden")
	ErrInvalid      = errors.New("invalid request")
	ErrExpired      = errors.New("session expired")
	ErrUnauthorized = errors.New("unauthorized")
	ErrUnavailable  = errors.New("dependency unavailable")
)

type Role string

const (
	RoleCounselor Role = "counselor"
	RoleReviewer  Role = "reviewer"
)

type RecommendationStatus string

const (
	RecommendationDraft     RecommendationStatus = "draft"
	RecommendationSubmitted RecommendationStatus = "submitted"
	RecommendationApproved  RecommendationStatus = "approved"
	RecommendationArchived  RecommendationStatus = "archived"
)

type FollowupStatus string

const (
	FollowupPending FollowupStatus = "pending"
	FollowupRunning FollowupStatus = "running"
	FollowupDone    FollowupStatus = "done"
	FollowupFailed  FollowupStatus = "failed"
)

func (s RecommendationStatus) CanMove(to RecommendationStatus) bool {
	switch s {
	case RecommendationDraft:
		return to == RecommendationSubmitted
	case RecommendationSubmitted:
		return to == RecommendationApproved
	case RecommendationApproved:
		return to == RecommendationArchived
	default:
		return false
	}
}

type User struct {
	ID, Email, Name string
	Role            Role
}
type Student struct {
	ID, OwnerID, Name, Region string
	GraduationYear, Version   int
}
type Major struct{ ID, Code, Name, Field string }
type University struct{ ID, Name, Province string }
type AdmissionPlan struct {
	ID, UniversityID, MajorID string
	Year, Quota, Remaining    int
}
type CareerPath struct{ ID, MajorID, Title, Outlook string }
type Recommendation struct {
	ID, StudentID, CreatedBy, Note string
	Status                         RecommendationStatus
	Version                        int
}
type RecommendationItem struct {
	ID, RecommendationID, AdmissionPlanID, CareerPathID string
	Priority                                            int
}
type Followup struct {
	ID, RecommendationID string
	DueAt                string
	Status               FollowupStatus
	Attempts             int
}
type AuditEvent struct{ ID, ActorID, EntityType, EntityID, Action, Result, RequestID, Detail, CreatedAt string }

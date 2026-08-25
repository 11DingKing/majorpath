package domain

import (
	"fmt"
	"strings"
	"time"
)

type StudentProfile struct {
	Name           string
	GraduationYear int
	Region         string
	Interests      map[string]float64
}

func (p StudentProfile) Validate(now time.Time) error {
	if strings.TrimSpace(p.Name) == "" {
		return fmt.Errorf("%w: name", ErrInvalid)
	}
	if p.GraduationYear < now.Year() || p.GraduationYear > now.Year()+8 {
		return fmt.Errorf("%w: graduation year", ErrInvalid)
	}
	if strings.TrimSpace(p.Region) == "" {
		return fmt.Errorf("%w: region", ErrInvalid)
	}
	if len(p.Interests) > 12 {
		return fmt.Errorf("%w: too many interests", ErrInvalid)
	}
	for tag, w := range p.Interests {
		if strings.TrimSpace(tag) == "" || w <= 0 || w > 1 {
			return fmt.Errorf("%w: interest weight", ErrInvalid)
		}
	}
	return nil
}
func ValidateTransition(from, to RecommendationStatus) error {
	if from == to {
		return fmt.Errorf("%w: unchanged state", ErrConflict)
	}
	if !from.CanMove(to) {
		return fmt.Errorf("%w: %s to %s", ErrConflict, from, to)
	}
	return nil
}
func ValidateFollowup(s FollowupStatus) error {
	switch s {
	case FollowupPending, FollowupRunning, FollowupDone, FollowupFailed:
		return nil
	default:
		return fmt.Errorf("%w: followup status", ErrInvalid)
	}
}

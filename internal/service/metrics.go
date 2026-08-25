package service

import (
	"context"
	"fmt"
	"github.com/11DingKing/majorpath/internal/domain"
)

type Metrics struct {
	Students        int
	Recommendations int
	Approved        int
	Followups       int
}

func (m Metrics) Validate() error {
	if m.Students < 0 {
		return fmt.Errorf("%w: students", domain.ErrInvalid)
	}
	if m.Recommendations < 0 {
		return fmt.Errorf("%w: recommendations", domain.ErrInvalid)
	}
	if m.Approved < 0 || m.Approved > m.Recommendations {
		return fmt.Errorf("%w: approved", domain.ErrInvalid)
	}
	if m.Followups < 0 {
		return fmt.Errorf("%w: followups", domain.ErrInvalid)
	}
	return nil
}

func (m Metrics) ApprovalRate() float64 {
	if m.Recommendations == 0 {
		return 0
	}
	return float64(m.Approved) / float64(m.Recommendations)
}

func (m Metrics) Merge(other Metrics) Metrics {
	return Metrics{Students: m.Students + other.Students, Recommendations: m.Recommendations + other.Recommendations, Approved: m.Approved + other.Approved, Followups: m.Followups + other.Followups}
}

func (m Metrics) CheckContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return m.Validate()
	}
}

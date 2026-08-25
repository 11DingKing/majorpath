package auth

import (
	"context"
	"fmt"
	"github.com/11DingKing/majorpath/internal/domain"
)

type Policy struct{}

func (Policy) CanCreateStudent(u domain.User) error {
	if u.Role != domain.RoleCounselor {
		return fmt.Errorf("%w: counselor required", domain.ErrForbidden)
	}
	return nil
}
func (Policy) CanSubmit(u domain.User, owner string) error {
	if u.Role != domain.RoleCounselor || u.ID != owner {
		return domain.ErrForbidden
	}
	return nil
}
func (Policy) CanReview(u domain.User) error {
	if u.Role != domain.RoleReviewer {
		return domain.ErrForbidden
	}
	return nil
}
func (Policy) ContextUser(ctx context.Context) (domain.User, error) {
	u, ok := ctx.Value(userContextKey{}).(domain.User)
	if !ok {
		return domain.User{}, domain.ErrUnauthorized
	}
	return u, nil
}

type userContextKey struct{}

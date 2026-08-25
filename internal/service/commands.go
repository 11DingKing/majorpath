package service

import (
	"context"
	"fmt"
	"github.com/11DingKing/majorpath/internal/domain"
)

type CommandLog interface {
	Record(context.Context, string, string) error
}
type NoopLog struct{}

func (NoopLog) Record(context.Context, string, string) error { return nil }
func ValidateCommand(u domain.User, expected domain.Role, entity string) error {
	if u.ID == "" {
		return domain.ErrUnauthorized
	}
	if u.Role != expected {
		return fmt.Errorf("%w: expected %s", domain.ErrForbidden, expected)
	}
	if entity == "" {
		return fmt.Errorf("%w: entity required", domain.ErrInvalid)
	}
	return nil
}
func Execute(ctx context.Context, log CommandLog, action, entity string) error {
	if err := EnsureContext(ctx); err != nil {
		return err
	}
	if err := log.Record(ctx, action, entity); err != nil {
		return fmt.Errorf("record command: %w", err)
	}
	return nil
}

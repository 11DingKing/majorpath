package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/11DingKing/majorpath/internal/clock"
	"github.com/11DingKing/majorpath/internal/domain"
	"github.com/11DingKing/majorpath/internal/repository"
	"time"
)

type Service struct {
	Users    repository.UserRepo
	Sessions repository.SessionRepo
	Clock    clock.Clock
	TTL      time.Duration
}

func hash(v string) string { h := sha256.Sum256([]byte(v)); return hex.EncodeToString(h[:]) }
func (s Service) Login(ctx context.Context, email, password string) (string, domain.User, error) {
	u, p, e := s.Users.ByEmail(ctx, email)
	if e != nil {
		return "", u, domain.ErrUnauthorized
	}
	if p != hash(password) {
		return "", u, domain.ErrUnauthorized
	}
	raw := fmt.Sprintf("%s-%d", u.ID, s.Clock.Now().UnixNano())
	if e = s.Sessions.Create(ctx, fmt.Sprintf("%x", hash(raw)), u.ID, hash(raw), s.Clock.Now().Add(s.TTL).Format("2006-01-02 15:04:05")); e != nil {
		return "", u, e
	}
	return raw, u, nil
}
func (s Service) Authenticate(ctx context.Context, token string) (domain.User, error) {
	id, e := s.Sessions.Find(ctx, hash(token))
	if e != nil {
		return domain.User{}, domain.ErrUnauthorized
	}
	return s.Users.ByID(ctx, id)
}
func (s Service) Logout(ctx context.Context, token string) error {
	return s.Sessions.Revoke(ctx, hash(token))
}

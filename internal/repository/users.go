package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/11DingKing/majorpath/internal/domain"
)

type UserRepo struct{ DB *sql.DB }

func (r UserRepo) Create(ctx context.Context, u domain.User, password string) error {
	_, e := r.DB.ExecContext(ctx, "INSERT INTO users(id,email,name,role,password_hash,created_at) VALUES(?,?,?,?,?,datetime('now'))", u.ID, u.Email, u.Name, u.Role, password)
	return wrap("create user", e)
}
func (r UserRepo) ByEmail(ctx context.Context, email string) (domain.User, string, error) {
	var u domain.User
	var p string
	e := r.DB.QueryRowContext(ctx, "SELECT id,email,name,role,password_hash FROM users WHERE email=?", email).Scan(&u.ID, &u.Email, &u.Name, &u.Role, &p)
	return u, p, wrap("find user", e)
}
func (r UserRepo) ByID(ctx context.Context, id string) (domain.User, error) {
	var u domain.User
	e := r.DB.QueryRowContext(ctx, "SELECT id,email,name,role FROM users WHERE id=?", id).Scan(&u.ID, &u.Email, &u.Name, &u.Role)
	return u, wrap("find user", e)
}
func requireRole(u domain.User, role domain.Role) error {
	if u.Role != role {
		return fmt.Errorf("%w: role %s required", domain.ErrForbidden, role)
	}
	return nil
}

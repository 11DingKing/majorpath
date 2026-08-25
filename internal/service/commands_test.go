package service

import (
	"context"
	"errors"
	"github.com/11DingKing/majorpath/internal/domain"
	"testing"
)

type commandLog struct {
	actions []string
	err     error
}

func (l *commandLog) Record(_ context.Context, a, e string) error {
	if l.err != nil {
		return l.err
	}
	l.actions = append(l.actions, a+":"+e)
	return nil
}
func TestValidateCommand(t *testing.T) {
	good := domain.User{ID: "u", Role: domain.RoleCounselor}
	if e := ValidateCommand(good, domain.RoleCounselor, "student"); e != nil {
		t.Fatal(e)
	}
	cases := []struct {
		u      domain.User
		r      domain.Role
		entity string
	}{{domain.User{}, domain.RoleCounselor, "x"}, {good, domain.RoleReviewer, "x"}, {good, domain.RoleCounselor, ""}}
	for _, c := range cases {
		if e := ValidateCommand(c.u, c.r, c.entity); e == nil {
			t.Fatal("invalid command accepted")
		}
	}
}
func TestExecuteRecords(t *testing.T) {
	l := &commandLog{}
	if e := Execute(context.Background(), l, "create", "student"); e != nil {
		t.Fatal(e)
	}
	if len(l.actions) != 1 || l.actions[0] != "create:student" {
		t.Fatal(l.actions)
	}
}
func TestExecutePropagatesLoggerError(t *testing.T) {
	l := &commandLog{err: errors.New("disk")}
	if e := Execute(context.Background(), l, "x", "y"); e == nil || !errors.Is(e, l.err) {
		t.Fatal(e)
	}
}
func TestExecuteCancelled(t *testing.T) {
	ctx, c := context.WithCancel(context.Background())
	c()
	l := &commandLog{}
	if e := Execute(ctx, l, "x", "y"); !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
	if len(l.actions) != 0 {
		t.Fatal(l.actions)
	}
}

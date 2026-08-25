package service

import (
	"context"
	"errors"
	"github.com/11DingKing/majorpath/internal/domain"
	"testing"
	"time"
)

func TestLifecycleDeadline(t *testing.T) {
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	l := Lifecycle{Clock: func() time.Time { return base }}
	if e := l.Deadline(base.Add(time.Hour)); e != nil {
		t.Fatal(e)
	}
	if e := l.Deadline(base.Add(-time.Hour)); !errors.Is(e, domain.ErrExpired) {
		t.Fatal(e)
	}
}
func TestLifecycleContext(t *testing.T) {
	ctx, c := context.WithCancel(context.Background())
	if e := (Lifecycle{}).RequireActive(ctx); e != nil {
		t.Fatal(e)
	}
	c()
	if e := (Lifecycle{}).RequireActive(ctx); !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
}
func TestFollowupLifecycle(t *testing.T) {
	l := Lifecycle{}
	cases := []struct {
		current domain.FollowupStatus
		success bool
		want    domain.FollowupStatus
	}{{domain.FollowupPending, false, domain.FollowupRunning}, {domain.FollowupRunning, true, domain.FollowupDone}, {domain.FollowupRunning, false, domain.FollowupFailed}, {domain.FollowupFailed, false, domain.FollowupRunning}}
	for _, c := range cases {
		got, e := l.NextFollowupStatus(c.current, c.success)
		if e != nil || got != c.want {
			t.Fatalf("%+v got %s %v", c, got, e)
		}
	}
	if _, e := l.NextFollowupStatus(domain.FollowupDone, true); e == nil {
		t.Fatal("done transitioned")
	}
}
func TestRetryAllowed(t *testing.T) {
	l := Lifecycle{}
	for _, c := range []struct {
		a, m int
		ok   bool
	}{{0, 3, true}, {1, 3, true}, {2, 3, true}, {3, 3, false}, {-1, 3, false}, {1, 0, false}} {
		if got := l.RetryAllowed(c.a, c.m); got != c.ok {
			t.Fatalf("%+v got %v", c, got)
		}
	}
}
func TestLifecycleZeroClock(t *testing.T) {
	l := Lifecycle{}
	if e := l.RequireActive(context.Background()); e != nil {
		t.Fatal(e)
	}
}

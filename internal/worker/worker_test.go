package worker

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetryDelay(t *testing.T) {
	p := RetryPolicy{Base: time.Second}
	if p.Delay(1) != time.Second || p.Delay(2) != 2*time.Second || p.Delay(4) != 8*time.Second {
		t.Fatal("delay")
	}
	if p.Delay(20) != 512*time.Second {
		t.Fatal("cap")
	}
}
func TestRunWithRetryEventuallySucceeds(t *testing.T) {
	n := 0
	e := RunWithRetry(context.Background(), RetryPolicy{Max: 3, Base: time.Millisecond}, func(context.Context) error {
		n++
		if n < 3 {
			return errors.New("temporary")
		}
		return nil
	})
	if e != nil || n != 3 {
		t.Fatalf("e=%v n=%d", e, n)
	}
}
func TestRunWithRetryExhausts(t *testing.T) {
	n := 0
	e := RunWithRetry(context.Background(), RetryPolicy{Max: 2, Base: time.Millisecond}, func(context.Context) error { n++; return errors.New("down") })
	if e == nil || n != 2 {
		t.Fatalf("e=%v n=%d", e, n)
	}
}
func TestRunWithRetryCancelled(t *testing.T) {
	ctx, c := context.WithCancel(context.Background())
	c()
	if e := RunWithRetry(ctx, RetryPolicy{Max: 2, Base: time.Millisecond}, func(context.Context) error { return errors.New("x") }); e == nil {
		t.Fatal("cancel ignored")
	}
}
func TestDelayMinimumAttempt(t *testing.T) {
	p := RetryPolicy{Base: 10 * time.Millisecond}
	if p.Delay(0) != 10*time.Millisecond {
		t.Fatal("minimum")
	}
}

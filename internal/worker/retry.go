package worker

import (
	"context"
	"fmt"
	"time"
)

type RetryPolicy struct {
	Max  int
	Base time.Duration
}

func (p RetryPolicy) Delay(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	if attempt > 10 {
		attempt = 10
	}
	d := p.Base
	for i := 1; i < attempt; i++ {
		d *= 2
	}
	return d
}
func RunWithRetry(ctx context.Context, p RetryPolicy, fn func(context.Context) error) error {
	var last error
	for i := 1; i <= p.Max; i++ {
		if e := fn(ctx); e == nil {
			return nil
		} else {
			last = e
		}
		t := time.NewTimer(p.Delay(i))
		select {
		case <-ctx.Done():
			t.Stop()
			return ctx.Err()
		case <-t.C:
		}
	}
	return fmt.Errorf("retry exhausted: %w", last)
}

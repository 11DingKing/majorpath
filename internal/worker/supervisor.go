package worker

import (
	"context"
	"sync"
	"time"
)

type Supervisor struct {
	Workers []func(context.Context)
	Wait    time.Duration
}

func (s Supervisor) Run(ctx context.Context) {
	var wg sync.WaitGroup
	for _, fn := range s.Workers {
		if fn == nil {
			continue
		}
		wg.Add(1)
		go func(f func(context.Context)) { defer wg.Done(); f(ctx) }(fn)
	}
	if s.Wait > 0 {
		timer := time.NewTimer(s.Wait)
		select {
		case <-ctx.Done():
			timer.Stop()
		case <-timer.C:
		}
	}
	wg.Wait()
}

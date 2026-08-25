package worker

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestSupervisorRunsWorkers(t *testing.T) {
	var n atomic.Int32
	s := Supervisor{Workers: []func(context.Context){func(context.Context) { n.Add(1) }, func(context.Context) { n.Add(1) }, nil}, Wait: time.Millisecond}
	s.Run(context.Background())
	if n.Load() != 2 {
		t.Fatalf("n=%d", n.Load())
	}
}
func TestSupervisorStopsOnContext(t *testing.T) {
	ctx, c := context.WithCancel(context.Background())
	started := make(chan struct{})
	s := Supervisor{Workers: []func(context.Context){func(ctx context.Context) { close(started); <-ctx.Done() }}, Wait: time.Second}
	done := make(chan struct{})
	go func() { s.Run(ctx); close(done) }()
	<-started
	c()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("did not stop")
	}
}
func TestSupervisorEmpty(t *testing.T) { Supervisor{}.Run(context.Background()) }

package config

import (
	"os"
	"testing"
	"time"
)

func TestDefaults(t *testing.T) {
	for _, k := range []string{"PORT", "DATABASE_PATH", "SESSION_TTL", "WORKER_INTERVAL"} {
		os.Unsetenv(k)
	}
	c := Load()
	if c.Port != "8080" || c.DatabasePath != "./majorpath.db" {
		t.Fatalf("defaults: %+v", c)
	}
	if c.SessionTTL != 12*time.Hour || c.WorkerInterval != 5*time.Second {
		t.Fatalf("durations: %+v", c)
	}
}
func TestEnvironment(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("DATABASE_PATH", "/tmp/x.db")
	t.Setenv("SESSION_TTL", "30m")
	t.Setenv("WORKER_INTERVAL", "2")
	c := Load()
	if c.Port != "9090" || c.DatabasePath != "/tmp/x.db" {
		t.Fatal(c)
	}
	if c.SessionTTL != 30*time.Minute || c.WorkerInterval != 2*time.Second {
		t.Fatal(c)
	}
}
func TestInvalidDurationUsesDefault(t *testing.T) {
	t.Setenv("SESSION_TTL", "bad")
	c := Load()
	if c.SessionTTL != 12*time.Hour {
		t.Fatalf("ttl=%v", c.SessionTTL)
	}
}

package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port           string
	DatabasePath   string
	SessionTTL     time.Duration
	WorkerInterval time.Duration
}

func Load() Config {
	c := Config{Port: env("PORT", "8080"), DatabasePath: env("DATABASE_PATH", "./majorpath.db"), SessionTTL: duration("SESSION_TTL", 12*time.Hour), WorkerInterval: duration("WORKER_INTERVAL", 5*time.Second)}
	return c
}
func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
func duration(k string, d time.Duration) time.Duration {
	if v := os.Getenv(k); v != "" {
		if x, e := time.ParseDuration(v); e == nil {
			return x
		}
		if n, e := strconv.Atoi(v); e == nil {
			return time.Duration(n) * time.Second
		}
	}
	return d
}

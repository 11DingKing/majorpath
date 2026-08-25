package clock

import (
	"testing"
	"time"
)

func TestFixedClock(t *testing.T) {
	want := time.Date(2026, 8, 25, 3, 4, 5, 0, time.UTC)
	if got := (Fixed{T: want}).Now(); !got.Equal(want) {
		t.Fatal(got)
	}
}
func TestRealClockUTC(t *testing.T) {
	before := time.Now().UTC()
	got := (Real{}).Now()
	after := time.Now().UTC()
	if got.Before(before) || got.After(after) {
		t.Fatalf("outside range %v", got)
	}
	if got.Location() != time.UTC {
		t.Fatal(got.Location())
	}
}
func TestClockInterface(t *testing.T) {
	var c Clock = Fixed{T: time.Unix(0, 0).UTC()}
	if c.Now().Unix() != 0 {
		t.Fatal(c.Now())
	}
}

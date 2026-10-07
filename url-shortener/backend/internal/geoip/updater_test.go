package geoip

import (
	"testing"
	"time"
)

func TestNextUpdate(t *testing.T) {
	for _, sample := range []struct{ now, want string }{
		{"2026-10-08T12:00:00+03:00", "2026-11-01T00:00:00+03:00"},
		{"2026-12-31T23:59:59+03:00", "2027-01-01T00:00:00+03:00"},
		{"2026-11-01T00:00:00+03:00", "2026-12-01T00:00:00+03:00"},
	} {
		now, _ := time.Parse(time.RFC3339, sample.now)
		if got := nextUpdate(now).Format(time.RFC3339); got != sample.want {
			t.Fatalf("%s: got %s, want %s", sample.now, got, sample.want)
		}
	}
}

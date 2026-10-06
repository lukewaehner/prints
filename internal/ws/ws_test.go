package ws

import (
	"testing"
	"time"
)

func TestJitter_delay_stays_within_half_to_full_backoff(t *testing.T) {
	tests := []struct {
		name    string
		attempt int
		min     time.Duration
		max     time.Duration
	}{
		{name: "first attempt uses base", attempt: 0, min: 500 * time.Millisecond, max: 1 * time.Second},
		{name: "second attempt doubles", attempt: 1, min: 1 * time.Second, max: 2 * time.Second},
		{name: "third attempt doubles again", attempt: 2, min: 2 * time.Second, max: 4 * time.Second},
		{name: "reaching cap is clamped to cap", attempt: 3, min: 4 * time.Second, max: 8 * time.Second},
		{name: "past cap stays at cap", attempt: 10, min: 4 * time.Second, max: 8 * time.Second},
		{name: "huge attempt does not overflow", attempt: 1000, min: 4 * time.Second, max: 8 * time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// act
			got := jitter(1*time.Second, 8*time.Second, tt.attempt)

			// assert
			if got < tt.min || got > tt.max {
				t.Errorf("jitter(attempt=%d) = %v, want between %v and %v", tt.attempt, got, tt.min, tt.max)
			}
		})
	}
}

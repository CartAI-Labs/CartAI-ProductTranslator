package infrastructure

import (
	"testing"
	"time"
)

func TestNextBackoff(t *testing.T) {
	tests := []struct {
		name    string
		attempt int
		want    time.Duration
	}{
		{"first failure", 1, 500 * time.Millisecond},
		{"second failure doubles", 2, 1 * time.Second},
		{"third failure doubles again", 3, 2 * time.Second},
		{"caps at max instead of growing forever", 10, 30 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := nextBackoff(tt.attempt)
			if got != tt.want {
				t.Errorf("nextBackoff(%d) = %v, want %v", tt.attempt, got, tt.want)
			}
		})
	}
}

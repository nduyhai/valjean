package limiter

import (
	"context"
	"time"

	"go.uber.org/ratelimit"
)

type Memory struct {
	rl      ratelimit.Limiter
	timeout time.Duration
}

func NewMemory() *Memory {
	rl := ratelimit.New(2) // request per second
	return &Memory{rl: rl, timeout: 100 * time.Millisecond}
}

func (m *Memory) Allow(ctx context.Context, key string, tokens int) (bool, int) {
	// Use goroutine with timeout to make it non-blocking
	done := make(chan struct{})
	go func() {
		m.rl.Take() // This blocks until token available
		close(done)
	}()

	select {
	case <-done:
		return true, 0 // Got token
	case <-time.After(m.timeout): // Quick timeout
		return false, 0 // Rate limited - DROP MESSAGE
	case <-ctx.Done():
		return false, 0 // Context cancelled
	}
}

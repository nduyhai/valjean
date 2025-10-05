package ports

import "context"

type RateLimiter interface {
	Allow(ctx context.Context, key string, tokens int) (bool, int) // ok, remaining
}

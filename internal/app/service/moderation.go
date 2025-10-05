package service

import "context"

type Moderation interface {
	Allowed(ctx context.Context, text string) (bool, string) // ok, reason
}

type moderation struct {
}

func NewModeration() Moderation {
	return &moderation{}
}

func (m *moderation) Allowed(ctx context.Context, text string) (bool, string) {
	return true, ""
}

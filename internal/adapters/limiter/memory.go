package limiter

import "context"

type Memory struct {
}

func NewMemory() *Memory {
	return &Memory{}
}

func (m *Memory) Allow(ctx context.Context, key string, tokens int) (bool, int) {
	return true, 1
}

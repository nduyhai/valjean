package service

import (
	"context"

	"github.com/nduyhai/valjean/internal/app/entities"
	"github.com/nduyhai/valjean/internal/infra/config"
)

const PrivateChat = "private"

type Moderation interface {
	Allowed(ctx context.Context, in entities.EvalInput) bool // ok, reason
}

type moderationStrategy interface {
	Allowed(ctx context.Context, in entities.EvalInput) bool
}

type moderation struct {
	strategies map[entities.SourceType]moderationStrategy
}

func NewModeration(config config.Config) Moderation {
	return &moderation{
		strategies: map[entities.SourceType]moderationStrategy{
			entities.SourceTelegram: newTelegramModeration(config.Telegram),
			entities.SourceZalo:     newZaloModeration(config.Zalo),
		},
	}
}

func (m *moderation) Allowed(ctx context.Context, in entities.EvalInput) bool {
	strategy, ok := m.strategies[in.SourceType]
	if !ok {
		return false
	}

	return strategy.Allowed(ctx, in)
}

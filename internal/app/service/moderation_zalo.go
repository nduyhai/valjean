package service

import (
	"context"

	"github.com/nduyhai/valjean/internal/app/entities"
	"github.com/nduyhai/valjean/internal/infra/config"
)

type zaloModeration struct {
	config config.Zalo
}

func newZaloModeration(config config.Zalo) moderationStrategy {
	return &zaloModeration{config: config}
}

func (m *zaloModeration) Allowed(_ context.Context, _ entities.EvalInput) bool {
	return true
}

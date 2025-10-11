package service

import (
	"context"
	"slices"
	"strings"

	"github.com/nduyhai/valjean/internal/app/entities"
	"github.com/nduyhai/valjean/internal/infra/config"
)

type zaloModeration struct {
	config config.Zalo
}

func newZaloModeration(config config.Zalo) moderationStrategy {
	return &zaloModeration{config: config}
}

func (m *zaloModeration) Allowed(_ context.Context, in entities.EvalInput) bool {
	text := strings.TrimSpace(in.Text)

	if len(m.config.BlockedUsers) > 0 && slices.Contains(m.config.BlockedUsers, in.UserHandle) {
		return false
	}

	if len(m.config.AllowedUsers) > 0 && !slices.Contains(m.config.AllowedUsers, in.UserHandle) {
		return false
	}

	if m.config.BotUsername != "" &&
		strings.Contains(text, "@"+m.config.BotUsername) {
		return true
	}

	if in.ReplyFor == m.config.BotUsername {
		return true
	}

	if in.ChatType == PrivateChat {
		return true
	}

	return false
}

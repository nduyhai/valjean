package service

import (
	"context"
	"slices"
	"strings"

	"github.com/nduyhai/valjean/internal/app/entities"
	"github.com/nduyhai/valjean/internal/infra/config"
)

type Moderation interface {
	Allowed(ctx context.Context, in entities.EvalInput) bool // ok, reason
}

type moderation struct {
	telegram config.Telegram
}

func NewModeration(config config.Config) Moderation {
	return &moderation{telegram: config.Telegram}
}

func (m *moderation) Allowed(ctx context.Context, in entities.EvalInput) bool {
	text := strings.TrimSpace(in.Text)

	if len(m.telegram.BlockedUsers) > 0 && slices.Contains(m.telegram.BlockedUsers, in.UserHandle) {
		return false
	}

	if len(m.telegram.AllowedUsers) > 0 && !slices.Contains(m.telegram.AllowedUsers, in.UserHandle) {
		return false
	}

	if m.telegram.Prefix != "" && strings.HasPrefix(text, m.telegram.Prefix) {
		return true
	}
	if m.telegram.BotUsername != "" &&
		strings.Contains(text, "@"+m.telegram.BotUsername) {
		return true
	}

	if in.ReplyFor == m.telegram.BotUsername {
		return true
	}

	if in.ChatType == "private" {
		return true
	}

	return false
}

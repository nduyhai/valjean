package producer

import (
	"context"
	"log/slog"
	"time"

	zalobotapi "github.com/nduyhai/go-zalo-bot-api"
	"github.com/nduyhai/go-zalo-bot-api/endpoints"
	"github.com/nduyhai/valjean/internal/app/entities"
)

type Zalo struct {
	logger *slog.Logger
	client *zalobotapi.Client
}

func NewZalo(logger *slog.Logger, client *zalobotapi.Client) *Zalo {
	return &Zalo{logger: logger, client: client}
}

func (z *Zalo) Publish(ctx context.Context, event entities.Event) {
	resp, err := endpoints.SendMessage(ctx, z.client, endpoints.SendMessageReq{ChatID: event.ChatID, Text: event.ReplyMessage}, zalobotapi.WithTimeout(5*time.Second))
	if err != nil {
		z.logger.Error("sendMessage error:", slog.Any("error", err))
		return
	}
	z.logger.Info("sendMessage response:", slog.Any("response", resp))
}

func (z *Zalo) Supported() entities.SourceType {
	return entities.SourceZalo
}

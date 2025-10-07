package producer

import (
	"context"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nduyhai/valjean/internal/app/entities"
)

type Telegram struct {
	bot    *tgbotapi.BotAPI
	logger *slog.Logger
}

func NewTelegram(bot *tgbotapi.BotAPI, logger *slog.Logger) *Telegram {
	return &Telegram{bot: bot, logger: logger}
}

func (t *Telegram) Publish(ctx context.Context, event entities.Event) {

	message := tgbotapi.NewMessage(event.ChatID, event.ReplyMessage)
	message.ReplyToMessageID = event.OriginalMessageId

	_, err := t.bot.Send(message)
	if err != nil {
		t.logger.Error("failed to send message: ", slog.Any("error", err))
	}

}

package producer

import (
	"context"
	"log/slog"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nduyhai/valjean/internal/app/entities"
)

type Telegram struct {
	bot    *tgbotapi.BotAPI
	logger *slog.Logger
}

func (t *Telegram) Supported() entities.SourceType {
	return entities.SourceTelegram
}

func NewTelegram(bot *tgbotapi.BotAPI, logger *slog.Logger) *Telegram {
	return &Telegram{bot: bot, logger: logger}
}

func (t *Telegram) Publish(ctx context.Context, event entities.Event) {
	messageId, _ := strconv.Atoi(event.OriginalMessageId)
	chatId, _ := strconv.ParseInt(event.ChatID, 10, 64)

	safe := tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, event.ReplyMessage)
	message := tgbotapi.NewMessage(chatId, safe)
	message.ReplyToMessageID = messageId
	message.ParseMode = tgbotapi.ModeMarkdownV2

	_, err := t.bot.Send(message)
	if err != nil {
		t.logger.Error("failed to send message: ", slog.Any("error", err))
	}

}

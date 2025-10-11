package http

import (
	"log/slog"
	"net/http"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nduyhai/valjean/internal/infra/config"

	"github.com/gin-gonic/gin"
	"github.com/nduyhai/valjean/internal/app/entities"
	"github.com/nduyhai/valjean/internal/app/usecase"
)

type TelegramHandler struct {
	evaluator *usecase.EvaluateUseCase
	logger    *slog.Logger
	config    config.Config
}

func NewTelegramHandler(evaluator *usecase.EvaluateUseCase, logger *slog.Logger, config config.Config) *TelegramHandler {
	return &TelegramHandler{
		evaluator: evaluator,
		logger:    logger,
		config:    config,
	}
}

func (h *TelegramHandler) WebHook(c *gin.Context) {
	if !h.validateToken(c) {
		return
	}

	upd, ok := h.parseAndValidateUpdate(c)
	if !ok {
		return
	}

	h.logger.Info("received message:", upd.Message.Text, "")

	evalInput := entities.EvalInput{
		SourceType:   entities.SourceTelegram,
		ChatID:       strconv.FormatInt(upd.Message.Chat.ID, 10),
		MessageID:    strconv.Itoa(upd.Message.MessageID),
		UserID:       strconv.FormatInt(upd.Message.From.ID, 10),
		UserHandle:   upd.Message.From.UserName,
		Text:         upd.Message.Text,
		ContextSnips: h.extractContextSnips(upd),
		ChatType:     upd.Message.Chat.Type,
		ReplyFor:     h.getReplyUserName(upd),
	}

	err := h.evaluator.Handle(c.Request.Context(), evalInput)
	if err != nil {
		h.logger.Error("failed to handle message: ", slog.Any("error", err))
		return
	}
}

func (h *TelegramHandler) validateToken(c *gin.Context) bool {
	token := c.Param("token")
	if token != h.config.Telegram.WebhookSecret {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		h.logger.Error("invalid token")
		return false
	}
	return true
}

func (h *TelegramHandler) parseAndValidateUpdate(c *gin.Context) (tgbotapi.Update, bool) {
	var upd tgbotapi.Update
	if err := c.ShouldBindJSON(&upd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		h.logger.Warn("invalid payload")
		return upd, false
	}
	if upd.Message == nil || upd.Message.Chat == nil {
		return upd, false
	}
	return upd, true
}

func (h *TelegramHandler) getReplyUserName(upd tgbotapi.Update) string {
	replyUserName := ""
	if upd.Message.ReplyToMessage != nil &&
		upd.Message.ReplyToMessage.From != nil {
		replyUserName = upd.Message.ReplyToMessage.From.UserName
	}
	return replyUserName
}

func (h *TelegramHandler) extractContextSnips(upd tgbotapi.Update) []string {
	if upd.Message.ReplyToMessage != nil && upd.Message.ReplyToMessage.Text != "" {
		return []string{upd.Message.ReplyToMessage.Text}
	}
	return nil
}

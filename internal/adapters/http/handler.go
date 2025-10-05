package http

import (
	"log/slog"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nduyhai/valjean/internal/infra/config"

	"github.com/gin-gonic/gin"
	"github.com/nduyhai/valjean/internal/app/entities"
	"github.com/nduyhai/valjean/internal/app/usecase"
)

type Handler struct {
	evaluator *usecase.EvaluateUseCase
	logger    *slog.Logger
	config    config.Config
}

func NewHandler(evaluator *usecase.EvaluateUseCase, logger *slog.Logger, config config.Config) *Handler {
	return &Handler{
		evaluator: evaluator,
		logger:    logger,
		config:    config,
	}
}

func (h *Handler) WebHook(c *gin.Context) {
	token := c.Param("token")
	if token != h.config.Telegram.WebhookSecret {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		h.logger.Error("invalid token")
		return
	}

	var upd tgbotapi.Update
	if err := c.ShouldBindJSON(&upd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		h.logger.Warn("invalid payload")
		return
	}

	if upd.Message == nil || upd.Message.Chat == nil {
		return
	}
	h.logger.Info("received message:", upd.Message.Text, "")

	var contextSnips []string
	if upd.Message.ReplyToMessage != nil && upd.Message.ReplyToMessage.Text != "" {
		contextSnips = []string{upd.Message.ReplyToMessage.Text}
	}

	evalInput := entities.EvalInput{
		ChatID:       upd.Message.Chat.ID,
		MessageID:    upd.Message.MessageID,
		UserID:       upd.Message.From.ID,
		UserHandle:   upd.Message.From.UserName,
		Text:         upd.Message.Text,
		ContextSnips: contextSnips,
		ChatType:     upd.Message.Chat.Type,
		ReplyFor:     h.getReplyUserName(upd),
	}

	shouldHandle := h.evaluator.ShouldHandle(c.Request.Context(), evalInput)

	if !shouldHandle {
		return
	}

	err := h.evaluator.Handle(c.Request.Context(), evalInput)
	if err != nil {
		return
	}
}

func (h *Handler) getReplyUserName(upd tgbotapi.Update) string {
	replyUserName := ""
	if upd.Message.ReplyToMessage != nil &&
		upd.Message.ReplyToMessage.From != nil {
		replyUserName = upd.Message.ReplyToMessage.From.UserName
	}
	return replyUserName
}

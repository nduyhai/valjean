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
	return &Handler{evaluator: evaluator, logger: logger, config: config}
}

func (h *Handler) WebHook(c *gin.Context) {
	token := c.Param("token")
	if token != h.config.Telegram.WebhookSecret {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		h.logger.Error("invalid token", nil, "")
		return
	}

	var upd tgbotapi.Update
	if err := c.ShouldBindJSON(&upd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		h.logger.Error("invalid payload", err, "")
		return
	}

	h.logger.Info("received message:", upd.Message.Text, "")
	err := h.evaluator.Handle(c.Request.Context(), entities.EvalInput{
		ChatID:       upd.Message.Chat.ID,
		MessageID:    upd.Message.MessageID,
		UserID:       upd.Message.From.ID,
		UserHandle:   upd.Message.From.UserName,
		Text:         upd.Message.Text,
		ContextSnips: nil,
	})
	if err != nil {
		return
	}
}

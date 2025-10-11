package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nduyhai/go-zalo-bot-api/endpoints"

	"github.com/nduyhai/valjean/internal/app/entities"
	"github.com/nduyhai/valjean/internal/app/usecase"
	"github.com/nduyhai/valjean/internal/infra/config"
)

type ZaloHandler struct {
	evaluator *usecase.EvaluateUseCase
	logger    *slog.Logger
	config    config.Config
}

func NewZaloHandler(evaluator *usecase.EvaluateUseCase, logger *slog.Logger, config config.Config) *ZaloHandler {
	return &ZaloHandler{evaluator: evaluator, logger: logger, config: config}
}

func (h *ZaloHandler) WebHook(c *gin.Context) {
	if !h.validateToken(c) {
		return
	}

	upd, ok := h.parseAndValidateUpdate(c)
	if !ok {
		return
	}

	h.logger.Info("received message:", upd.Message.Text, "")

	evalInput := entities.EvalInput{
		SourceType:   entities.SourceZalo,
		ChatID:       upd.Message.Chat.ID,
		MessageID:    upd.Message.MessageID,
		UserID:       upd.Message.From.ID,
		UserHandle:   upd.Message.From.DisplayName,
		Text:         h.getMessageText(upd),
		ContextSnips: h.extractContextSnips(upd),
		ChatType:     upd.Message.Chat.ChatType,
		ReplyFor:     h.getReplyUserName(upd),
	}

	err := h.evaluator.Handle(c.Request.Context(), evalInput)
	if err != nil {
		h.logger.Error("failed to handle message: ", slog.Any("error", err))
		return
	}
}

func (h *ZaloHandler) getMessageText(upd endpoints.UpdateResult) string {
	if upd.EventName == endpoints.EventMessageTextReceived {
		return upd.Message.Text
	}
	return ""
}

func (h *ZaloHandler) validateToken(c *gin.Context) bool {
	token := c.GetHeader(string(endpoints.HeaderSecretToken))
	if token != h.config.Zalo.WebhookSecret {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		h.logger.Error("invalid token")
		return false
	}
	return true
}

func (h *ZaloHandler) parseAndValidateUpdate(c *gin.Context) (endpoints.UpdateResult, bool) {
	var upd endpoints.UpdateResult
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

func (h *ZaloHandler) getReplyUserName(upd endpoints.UpdateResult) string {
	//TODO: not support yet
	return ""
}

func (h *ZaloHandler) extractContextSnips(upd endpoints.UpdateResult) []string {
	//TODO: not support yet
	return nil
}

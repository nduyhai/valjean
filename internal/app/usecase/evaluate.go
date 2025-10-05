package usecase

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/nduyhai/valjean/internal/app/entities"
	"github.com/nduyhai/valjean/internal/app/service"
	"github.com/nduyhai/valjean/internal/infra/config"
	"github.com/nduyhai/valjean/internal/ports"
)

type EvaluateUseCase struct {
	Evaluator     ports.Evaluator
	Moderation    service.Moderation
	RateLimiter   ports.RateLimiter
	EventProducer ports.EventProducer
	Telegram      config.Telegram
	Cooldown      time.Duration
	Logger        *slog.Logger
}

func NewEvaluateUseCase(evaluator ports.Evaluator, moderation service.Moderation, rateLimiter ports.RateLimiter, eventProducer ports.EventProducer, config config.Config, logger *slog.Logger) *EvaluateUseCase {
	return &EvaluateUseCase{Evaluator: evaluator, Moderation: moderation, RateLimiter: rateLimiter, EventProducer: eventProducer, Telegram: config.Telegram, Cooldown: 2 * time.Second, Logger: logger}

}

func (e *EvaluateUseCase) ShouldHandle(ctx context.Context, in entities.EvalInput) bool {
	text := strings.TrimSpace(in.Text)
	if e.Telegram.Prefix != "" && strings.HasPrefix(text, e.Telegram.Prefix) {
		return true
	}
	if e.Telegram.BotUsername != "" &&
		strings.Contains(text, "@"+e.Telegram.BotUsername) {
		return true
	}

	if in.ReplyFor == e.Telegram.BotUsername {
		return true
	}

	if in.ChatType == "private" {
		return true
	}

	return false
}

func (e *EvaluateUseCase) Handle(ctx context.Context, in entities.EvalInput) error {
	// rate limit
	ok, _ := e.RateLimiter.Allow(ctx, rlKey(in), 1)
	if !ok {
		e.Logger.Warn("rate limit exceeded")
		return errors.New("cooling down—try again in a moment")
	}
	// moderation
	allowed, reason := e.Moderation.Allowed(ctx, in.Text)
	if !allowed {
		e.Logger.Warn("message skipped")
		return errors.New("Message skipped (" + reason + ")")
	}
	// normalize text (trim trigger)
	in.Text = stripTriggers(in.Text, e.Telegram)
	out, err := e.Evaluator.Evaluate(ctx, in)
	if err != nil || out.Summary == "" {
		e.Logger.Error("failed to evaluate", slog.Any("error", err))
		return errors.New("i couldn’t evaluate that right now")
	}

	e.EventProducer.Publish(ctx, entities.Event{
		ChatID:            in.ChatID,
		OriginalMessageId: in.MessageID,
		ReplyMessage:      out.Summary,
	})
	return nil
}

func stripTriggers(s string, t config.Telegram) string {
	if t.Prefix != "" && strings.HasPrefix(strings.TrimSpace(s), t.Prefix) {
		return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(s), t.Prefix))
	}
	if t.BotUsername != "" {
		return strings.ReplaceAll(s, "@"+t.BotUsername, "")
	}
	return s
}

func rlKey(in entities.EvalInput) string {
	return "rl:user:" + strconv.Itoa(int(in.UserID))
}

package usecase

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"time"

	"github.com/nduyhai/valjean/internal/app/entities"
	"github.com/nduyhai/valjean/internal/app/service"
	"github.com/nduyhai/valjean/internal/infra/config"
	"github.com/nduyhai/valjean/internal/ports"
)

type EvaluateUseCase struct {
	evaluator     ports.Evaluator
	moderation    service.Moderation
	rateLimiter   ports.RateLimiter
	eventProducer ports.EventProducer
	telegram      config.Telegram
	cooldown      time.Duration
	logger        *slog.Logger
}

func NewEvaluateUseCase(evaluator ports.Evaluator, moderation service.Moderation, rateLimiter ports.RateLimiter, eventProducer ports.EventProducer, config config.Config, logger *slog.Logger) *EvaluateUseCase {
	return &EvaluateUseCase{evaluator: evaluator, moderation: moderation, rateLimiter: rateLimiter, eventProducer: eventProducer, telegram: config.Telegram, cooldown: 2 * time.Second, logger: logger}

}

func (e *EvaluateUseCase) Handle(ctx context.Context, in entities.EvalInput) error {
	// rate limit
	ok, _ := e.rateLimiter.Allow(ctx, rlKey(in), 1)
	if !ok {
		e.logger.Warn("rate limit exceeded")
		return errors.New("cooling down—try again in a moment")
	}
	// moderation
	allowed := e.moderation.Allowed(ctx, in)
	if !allowed {
		e.logger.Warn("message skipped")
		return errors.New("message skipped")
	}
	out, err := e.evaluator.Evaluate(ctx, in)
	if err != nil || out.Summary == "" {
		e.logger.Error("failed to evaluate", slog.Any("error", err))
		return errors.New("i couldn’t evaluate that right now")
	}

	e.eventProducer.Publish(ctx, entities.Event{
		ChatID:            in.ChatID,
		OriginalMessageId: in.MessageID,
		ReplyMessage:      out.Summary,
	})
	return nil
}

func rlKey(in entities.EvalInput) string {
	return "rl:user:" + strconv.Itoa(int(in.UserID))
}

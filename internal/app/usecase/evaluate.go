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
	worker        ports.Worker
	telegram      config.Telegram
	cooldown      time.Duration
	logger        *slog.Logger
}

func NewEvaluateUseCase(evaluator ports.Evaluator, moderation service.Moderation, rateLimiter ports.RateLimiter, eventProducer ports.EventProducer, worker ports.Worker, config config.Config, logger *slog.Logger) *EvaluateUseCase {
	return &EvaluateUseCase{
		evaluator:     evaluator,
		moderation:    moderation,
		rateLimiter:   rateLimiter,
		eventProducer: eventProducer,
		worker:        worker,
		telegram:      config.Telegram,
		cooldown:      5 * time.Second,
		logger:        logger,
	}

}

func (e *EvaluateUseCase) Handle(ctx context.Context, in entities.EvalInput) error {
	// rate limit
	ok, _ := e.rateLimiter.Allow(ctx, rlKey(in), 1)
	if !ok {
		e.logger.Warn("rate limit exceeded")

		e.sendMsg(ctx, in, "cooling down—try again in a moment")

		return errors.New("cooling down—try again in a moment")
	}

	e.worker.Submit(func() {
		e.process(in)
	})

	return nil
}

func (e *EvaluateUseCase) process(in entities.EvalInput) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), e.cooldown)
	defer cancelFunc()

	// moderation
	allowed := e.moderation.Allowed(ctx, in)
	if !allowed {
		e.logger.Warn("message skipped")
		return
	}
	out, err := e.evaluator.Evaluate(ctx, in)
	if err != nil || out.Summary == "" {

		e.logger.Error("failed to evaluate", slog.Any("error", err))

		e.sendMsg(ctx, in, "i couldn’t evaluate that right now")

		return
	}

	e.sendMsg(ctx, in, out.Summary)
}

func (e *EvaluateUseCase) sendMsg(ctx context.Context, in entities.EvalInput, replyMsg string) {
	e.eventProducer.Publish(ctx, entities.Event{
		ChatID:            in.ChatID,
		OriginalMessageId: in.MessageID,
		ReplyMessage:      replyMsg,
	})
}

func rlKey(in entities.EvalInput) string {
	return "rl:user:" + strconv.Itoa(int(in.UserID))
}

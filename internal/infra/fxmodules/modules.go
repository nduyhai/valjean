package fxmodules

import (
	"log/slog"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nduyhai/valjean/internal/adapters/limiter"
	"github.com/nduyhai/valjean/internal/adapters/llm/openai"
	"github.com/nduyhai/valjean/internal/adapters/producer"
	"github.com/nduyhai/valjean/internal/app/service"
	"github.com/nduyhai/valjean/internal/app/usecase"
	"github.com/nduyhai/valjean/internal/infra/config"
	"github.com/nduyhai/valjean/internal/ports"
	"go.uber.org/fx"
)

var ConfigModule = fx.Module("config",
	fx.Provide(config.Load, NewLogger),
)

func NewLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo, // change to LevelDebug for verbose
	}))
	slog.SetDefault(logger)
	slog.Info("Starting valjean...")
	return logger
}

var UseCaseModule = fx.Module("usecases",
	fx.Provide(
		usecase.NewEvaluateUseCase,
	),
)

var LimiterModule = fx.Module("adapters",
	fx.Provide(
		fx.Annotate(limiter.NewMemory, fx.As(new(ports.RateLimiter))),
		fx.Annotate(openai.NewClient, fx.As(new(ports.Evaluator))),
		NewBot,
		fx.Annotate(producer.NewTelegram, fx.As(new(ports.EventProducer))),
	),
)

func NewBot(config config.Config) (*tgbotapi.BotAPI, error) {
	return tgbotapi.NewBotAPI(config.Telegram.Token)
}

var ApplicationModule = fx.Module("application",
	fx.Provide(
		service.NewModeration,
	),
)

var AllModules = fx.Options(
	ConfigModule,
	LimiterModule,
	UseCaseModule,
	ApplicationModule,
)

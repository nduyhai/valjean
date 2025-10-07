package fxmodules

import (
	"log/slog"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nduyhai/valjean/internal/adapters/limiter"
	"github.com/nduyhai/valjean/internal/adapters/llm/openai"
	"github.com/nduyhai/valjean/internal/adapters/producer"
	"github.com/nduyhai/valjean/internal/adapters/worker"
	"github.com/nduyhai/valjean/internal/app/service"
	"github.com/nduyhai/valjean/internal/app/usecase"
	"github.com/nduyhai/valjean/internal/infra/config"
	"github.com/nduyhai/valjean/internal/ports"
	goopenai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
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
		NewOpenAISdk,
		fx.Annotate(openai.NewClient, fx.As(new(ports.Evaluator))),
		NewTelegramSdk,
		fx.Annotate(producer.NewTelegram, fx.As(new(ports.EventProducer))),
		fx.Annotate(worker.NewMemory, fx.As(new(ports.Worker))),
	),
)

func NewTelegramSdk(config config.Config) (*tgbotapi.BotAPI, error) {
	return tgbotapi.NewBotAPI(config.Telegram.Token)
}

func NewOpenAISdk(config config.Config) goopenai.Client {
	ai := goopenai.NewClient(
		option.WithAPIKey(config.OpenAI.Key), // defaults to os.LookupEnv("OPENAI_API_KEY")
	)
	return ai
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

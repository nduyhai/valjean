package fxmodules

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/nduyhai/valjean/internal/adapters/http"
	"github.com/nduyhai/valjean/internal/infra/config"
	"github.com/nduyhai/valjean/internal/infra/httpserver"
	"github.com/nduyhai/valjean/internal/ports"
	"go.uber.org/fx"
)

var HandlerModule = fx.Module("handlers",
	fx.Provide(
		http.NewTelegramHandler,
		http.NewZaloHandler,
	),
)

type HandlerParams struct {
	fx.In
	Telegram *http.TelegramHandler
	Zalo     *http.ZaloHandler
}

func NewGinRouter(handlers HandlerParams, logger *slog.Logger) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(gin.LoggerWithFormatter(func(p gin.LogFormatterParams) string {
		logger.Info("gin",
			"method", p.Method,
			"path", p.Path,
			"status", p.StatusCode,
			"latency", p.Latency.String(),
			"client", p.ClientIP,
		)
		return ""
	}))
	router.GET("/healthz", func(c *gin.Context) { c.String(200, "ok") })
	router.GET("/readyz", func(c *gin.Context) { c.String(200, "ready") })

	router.POST("/tl/webhook/:token", handlers.Telegram.WebHook)
	router.POST("/zl/webhook", handlers.Zalo.WebHook)

	return router
}

func NewHTTPServer(cfg config.Config, router *gin.Engine) *httpserver.Server {
	return httpserver.New(cfg.HTTP.Port, router)
}

var HTTPServerModule = fx.Module("http-server",
	fx.Provide(
		NewGinRouter,
		NewHTTPServer,
	),
)
var ServerModule = fx.Options(
	HandlerModule,
	HTTPServerModule,
)

type ServerParams struct {
	fx.In
	HTTPServer *httpserver.Server
	Worker     ports.Worker
}

func ServerLifecycle(lc fx.Lifecycle, servers ServerParams) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Start HTTP server
			go func() {
				if err := servers.HTTPServer.Start(); err != nil {
					// Log error but don't stop the application
					// as this might be expected during shutdown
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Shutdown HTTP server
			if err := servers.HTTPServer.Shutdown(ctx); err != nil {
				// Log error
			}

			servers.Worker.Shutdown()

			return nil
		},
	})
}

// LifecycleModule provides server lifecycle management
var LifecycleModule = fx.Module("lifecycle",
	fx.Invoke(ServerLifecycle),
)

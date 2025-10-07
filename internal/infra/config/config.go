package config

import "github.com/nduyhai/xcore/config/envloader"

type Config struct {
	HTTP     HTTPConfig
	Telegram Telegram
	OpenAI   OpenAI
}

type HTTPConfig struct {
	Port string `env:"HTTP_PORT" envDefault:"8080"`
}

type Telegram struct {
	Prefix        string   `env:"TELEGRAM_PREFIX" envDefault:"!eval"` // "!eval "
	BotUsername   string   `env:"TELEGRAM_BOT_USERNAME" envDefault:"valjean"`
	Token         string   `env:"TELEGRAM_BOT_TOKEN"`
	WebhookSecret string   `env:"TELEGRAM_WEBHOOK_SECRET"`
	BlockedUsers  []string `env:"TELEGRAM_BLOCKED_USERS"`
	AllowedUsers  []string `env:"TELEGRAM_ALLOWED_USERS"`
}

type OpenAI struct {
	Key string `env:"OPENAI_KEY" envDefault:""`
}

func Load() Config {
	config := Config{}
	_ = envloader.Load(&config)

	return config
}

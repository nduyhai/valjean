# Valjean

[![Go](https://img.shields.io/badge/go-1.24+-blue)](https://go.dev/)
[![License](https://img.shields.io/github/license/nduyhai/valjean)](LICENSE)

A smart Telegram bot powered by OpenAI that provides AI-powered conversation and evaluation capabilities with context awareness, rate limiting, and content moderation.

## ✨ Features

- 🤖 **OpenAI Integration**: Uses GPT-4 for intelligent responses
- 💬 **Context Awareness**: Maintains conversation context through reply messages
- 🚦 **Rate Limiting**: Prevents spam and manages usage quotas
- 🛡️ **Content Moderation**: Built-in content filtering and safety checks
- 📱 **Flexible Triggers**: Responds to mentions, prefixes, or direct replies
- 🏗️ **Clean Architecture**: Uses dependency injection with Uber FX
- 🚀 **Webhook Support**: Efficient real-time message processing
- ⚡ **High Performance**: Concurrent message handling with proper synchronization

## 🚀 Getting Started

### Prerequisites

- Go 1.24 or higher
- Telegram Bot Token (from [@BotFather](https://t.me/botfather))
- OpenAI API Key

### 📦 Installation

```bash
# Clone the repository
git clone https://github.com/nduyhai/valjean.git
cd valjean

# Install dependencies
go mod tidy

# Build the project
go build ./cmd/bot
```

### ⚙️ Configuration

Create a `.env` file or set environment variables:

```env
# Telegram Configuration
TELEGRAM_BOT_TOKEN=your_bot_token_here
TELEGRAM_BOT_USERNAME=your_bot_username
TELEGRAM_WEBHOOK_SECRET=your_webhook_secret
TELEGRAM_REQUIRED_MENTION=true
TELEGRAM_PREFIX=!eval

# OpenAI Configuration
OPENAI_KEY=your_openai_api_key

# HTTP Server
HTTP_PORT=8080
```

### 🏃 Running the Bot

#### Local Development
```bash
make run
```

#### Using Docker
```bash
docker build -t valjean .
docker run -p 8080:8080 --env-file .env valjean
```

### 🌐 Webhook Setup

Deploy the service with HTTPS and configure the Telegram webhook:

```bash
# Set webhook URL
curl -X POST \
 "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/setWebhook?url=https://YOUR.DOMAIN/telegram/webhook/$WEBHOOK_SECRET"
```

## 🎯 Usage

The bot responds to messages in several ways:

1. **Direct Mention**: `@your_bot_username Hello!`
2. **Prefix Command**: `!eval What is the weather like?`
3. **Reply to Bot**: Reply to any bot message for continued conversation
4. **Private Messages**: All messages in private chats

### Example Conversations

```
User: @valjean What is artificial intelligence?
Bot: Artificial intelligence (AI) refers to the simulation of human intelligence...

User: !eval Explain quantum computing
Bot: Quantum computing is a revolutionary computing paradigm...
```

## 🏗️ Project Structure

```
valjean/
├── cmd/bot/                    # Application entry point
├── internal/
│   ├── adapters/              # External integrations
│   │   ├── http/              # HTTP handlers
│   │   ├── llm/openai/        # OpenAI client
│   │   ├── limiter/           # Rate limiting
│   │   └── producer/telegram/ # Telegram message producer
│   ├── app/                   # Business logic
│   │   ├── entities/          # Domain entities
│   │   ├── service/           # Domain services
│   │   └── usecase/           # Application use cases
│   ├── infra/                 # Infrastructure
│   │   ├── config/            # Configuration
│   │   ├── fxmodules/         # Dependency injection modules
│   │   └── httpserver/        # HTTP server setup
│   └── ports/                 # Interface definitions
├── docs/                      # Documentation
├── Dockerfile                 # Container configuration
└── Makefile                   # Build automation
```

## 🔧 Configuration Options

| Environment Variable | Default | Description |
|---------------------|---------|-------------|
| `TELEGRAM_BOT_TOKEN` | - | Your Telegram bot token |
| `TELEGRAM_BOT_USERNAME` | `valjean` | Bot username for mentions |
| `TELEGRAM_WEBHOOK_SECRET` | - | Secret for webhook security |
| `TELEGRAM_REQUIRED_MENTION` | `true` | Require @mention in groups |
| `TELEGRAM_PREFIX` | `!eval` | Command prefix trigger |
| `OPENAI_KEY` | - | OpenAI API key |
| `HTTP_PORT` | `8080` | HTTP server port |

## 🛠️ Development

### Setting

```shell
# set webhook
curl -s -X POST "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/setWebhook" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://YOUR.DOMAIN/telegram/webhook/SECRET",
    "secret_token": "YOUR_HEADER_SECRET"
  }'

# check
curl -s "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/getWebhookInfo" | jq

# delete
curl -s "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/deleteWebhook"


# global commands
curl -s -X POST "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/setMyCommands" \
  -H "Content-Type: application/json" \
  -d '{
    "commands": [
      {"command":"eval","description":"Analyze a replied message or text"},
      {"command":"help","description":"How to use this bot"}
    ]
  }'

# restrict to all group chats
curl -s -X POST "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/setMyCommands" \
  -H "Content-Type: application/json" \
  -d '{
    "commands":[{"command":"eval","description":"Analyze a replied message or text"}],
    "scope":{"type":"all_group_chats"}
  }'


# name
curl -s -X POST "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/setMyName" \
  -d "name=Valjean"

# short description (shown in profiles)
curl -s -X POST "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/setMyShortDescription" \
  --data-urlencode "short_description=Reply-and-tag me to analyze messages."

# long description (shown on /start)
curl -s -X POST "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/setMyDescription" \
  --data-urlencode "description=I analyze the message you replied to when you tag me in groups."


curl -s -X POST "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/setMyDefaultAdministratorRights" \
  -H "Content-Type: application/json" \
  -d '{"rights":{"can_delete_messages":true,"can_manage_topics":true},"for_channels":false}'


curl -s -X POST "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/setChatMenuButton" \
  -H "Content-Type: application/json" \
  -d '{"menu_button":{"type":"commands"}}'

```

Group Privacy Mode (ON/OFF) → must use @BotFather → Bot Settings → Group Privacy → Turn off
(Required if you want the bot to trigger on mentions like @YourBot.)

### Running Tests
```bash
make test
```

### Linting
```bash
make lint
```

### Building
```bash
make build
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Telegram Bot API](https://core.telegram.org/bots/api) for bot functionality
- [OpenAI](https://openai.com/) for AI capabilities
- [Uber FX](https://uber-go.github.io/fx/) for dependency injection


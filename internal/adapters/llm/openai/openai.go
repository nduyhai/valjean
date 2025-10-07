package openai

import (
	"context"
	"log/slog"
	"time"

	"github.com/nduyhai/valjean/internal/app/entities"
	goopenai "github.com/openai/openai-go/v3"
)

type Client struct {
	ai     goopenai.Client
	logger *slog.Logger
}

func NewClient(ai goopenai.Client, logger *slog.Logger) *Client {
	return &Client{ai: ai, logger: logger}
}

func (c *Client) Evaluate(ctx context.Context, in entities.EvalInput) (entities.EvalOutput, error) {
	ctxTimeout, cancelFunc := context.WithTimeout(ctx, 10*time.Second)
	defer cancelFunc()

	var messages []goopenai.ChatCompletionMessageParamUnion

	for _, contextText := range in.ContextSnips {
		if contextText != "" {
			messages = append(messages, goopenai.UserMessage(contextText))
		}
	}

	messages = append(messages, goopenai.UserMessage(in.Text))

	chatCompletion, err := c.ai.Chat.Completions.New(ctxTimeout, goopenai.ChatCompletionNewParams{
		Messages: messages,
		Model:    goopenai.ChatModelGPT4o,
	})

	msg := "Please try again"
	if err != nil {
		return entities.EvalOutput{
			Summary: msg,
		}, err
	}

	if length := len(chatCompletion.Choices); length > 0 {
		msg = chatCompletion.Choices[0].Message.Content
	}
	return entities.EvalOutput{
		Summary:    msg,
		Citations:  []string{"Citation"},
		Confidence: 80,
	}, nil
}

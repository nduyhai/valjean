package openai

import (
	"context"
	"log/slog"

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

	// Combine context snippets with the main message
	var fullMessage string
	if len(in.ContextSnips) > 0 {
		fullMessage = "Context:\n"
		for _, contextText := range in.ContextSnips {
			if contextText != "" {
				fullMessage += contextText + "\n"
			}
		}
		fullMessage += "\nUser message:\n" + in.Text
	} else {
		fullMessage = in.Text
	}

	messages := []goopenai.ChatCompletionMessageParamUnion{
		goopenai.UserMessage(fullMessage),
	}

	chatCompletion, err := c.ai.Chat.Completions.New(ctx, goopenai.ChatCompletionNewParams{
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

package openai

import (
	"context"
	"log/slog"
	"strings"

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

	// 1) Select & compress context
	ctxSnips := uniqueNonEmpty(in.ContextSnips)
	if len(ctxSnips) > 3 {
		ctxSnips = ctxSnips[len(ctxSnips)-3:] // keep last 3
	}
	for i := range ctxSnips {
		ctxSnips[i] = truncateRunes(strings.TrimSpace(ctxSnips[i]), 240)
	}

	// 2) Build compact prompt
	//   ctx lines as "- ..."; then Q: <message>
	var b strings.Builder
	if len(ctxSnips) > 0 {
		for _, s := range ctxSnips {
			if s != "" {
				b.WriteString("- ")
				b.WriteString(s)
				b.WriteByte('\n')
			}
		}
		b.WriteByte('\n')
	}
	b.WriteString("Q: ")
	b.WriteString(strings.TrimSpace(in.Text))
	userPrompt := b.String()

	// 3) Messages: short system + one user
	messages := []goopenai.ChatCompletionMessageParamUnion{
		goopenai.SystemMessage("You are concise. If uncertain, say so briefly."),
		goopenai.UserMessage(userPrompt),
	}

	// 4) Call with a small model & caps
	chatCompletion, err := c.ai.Chat.Completions.New(ctx, goopenai.ChatCompletionNewParams{
		Messages:    messages,
		Model:       goopenai.ChatModelGPT4oMini, // cheaper than 4o
		MaxTokens:   goopenai.Int(120),           // cap output tokens
		Temperature: goopenai.Float(0.2),
	})

	msg := "Please try again"
	if err != nil {
		return entities.EvalOutput{Summary: msg}, err
	}
	if len(chatCompletion.Choices) > 0 && chatCompletion.Choices[0].Message.Content != "" {
		msg = strings.TrimSpace(chatCompletion.Choices[0].Message.Content)
	}
	return entities.EvalOutput{
		Summary:    msg,
		Citations:  nil,
		Confidence: 80,
	}, nil
}

func uniqueNonEmpty(ss []string) []string {
	seen := make(map[string]struct{}, len(ss))
	out := make([]string, 0, len(ss))
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		key := strings.ToLower(s) // or hash/fingerprint
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, s)
	}
	return out
}

func truncateRunes(s string, max int) string {
	rs := []rune(s)
	if len(rs) <= max {
		return s
	}
	return string(rs[:max]) + "…"
}

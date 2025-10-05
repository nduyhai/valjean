package openai

import (
	"context"
	"net/http"

	"github.com/nduyhai/valjean/internal/app/entities"
)

type Client struct {
	HTTP  *http.Client
	Model string
	URL   string // e.g. https://api.openai.com/v1/chat/completions
}

func NewClient() *Client {
	return &Client{}

}

func (c *Client) Evaluate(ctx context.Context, in entities.EvalInput) (entities.EvalOutput, error) {

	return entities.EvalOutput{
		Summary:    "Summary",
		Citations:  []string{"Citation"},
		Confidence: 80,
	}, nil
}

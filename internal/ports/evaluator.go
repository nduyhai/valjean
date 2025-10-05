package ports

import (
	"context"

	"github.com/nduyhai/valjean/internal/app/entities"
)

type Evaluator interface {
	Evaluate(ctx context.Context, in entities.EvalInput) (entities.EvalOutput, error)
}

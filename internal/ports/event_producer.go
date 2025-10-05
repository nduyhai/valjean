package ports

import (
	"context"

	"github.com/nduyhai/valjean/internal/app/entities"
)

type EventProducer interface {
	Publish(ctx context.Context, event entities.Event)
}

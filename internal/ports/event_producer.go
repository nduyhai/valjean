package ports

import (
	"context"

	"github.com/nduyhai/valjean/internal/app/entities"
)

type ChannelEventProducer interface {
	Publish(ctx context.Context, event entities.Event)

	Supported() entities.SourceType
}

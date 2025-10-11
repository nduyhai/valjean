package service

import (
	"context"

	"github.com/nduyhai/valjean/internal/app/entities"
	"github.com/nduyhai/valjean/internal/ports"
)

type EventProducer interface {
	Publish(ctx context.Context, event entities.Event)
}

type eventProducer struct {
	channels map[entities.SourceType]ports.ChannelEventProducer
}

func NewEventProducer(channels []ports.ChannelEventProducer) EventProducer {
	m := make(map[entities.SourceType]ports.ChannelEventProducer, len(channels))
	for _, ch := range channels {
		m[ch.Supported()] = ch
	}
	return &eventProducer{
		channels: m,
	}
}

func (e *eventProducer) Publish(ctx context.Context, event entities.Event) {

	if producer, ok := e.channels[event.SourceType]; ok {
		producer.Publish(ctx, event)
	}
}

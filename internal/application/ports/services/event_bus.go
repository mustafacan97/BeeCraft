package eventBus

import (
	"context"
	"platform/internal/domain"
)

// Subscription represents an event subscription.
type Subscription struct {
	EventName string
	Handler   func(ctx context.Context, event domain.Event) error
}

// EventBus is the interface for publishing and subscribing to events.
type EventBus interface {
	Publish(ctx context.Context, event domain.Event) error
	Subscribe(subscriber, eventName string, handler func(ctx context.Context, event domain.Event) error) (string, error)
	Unsubscribe(subscriptionId string) error
	Close()
}

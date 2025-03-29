package eventBus

import (
	"context"
	domainEvents "platform/internal/domain/events"
)

// Subscription represents an event subscription.
type Subscription struct {
	EventName string
	Handler   func(ctx context.Context, event domainEvents.Event) error
}

// EventBus is the interface for publishing and subscribing to events.
type EventBus interface {
	Publish(ctx context.Context, event domainEvents.Event) error
	Subscribe(subscriber, eventName string, handler func(ctx context.Context, event domainEvents.Event) error) (string, error)
	Unsubscribe(subscriptionId string) error
	Close()
}

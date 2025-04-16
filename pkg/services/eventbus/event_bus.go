package eventbus

import (
	"context"
	"platform/pkg/domain"
)

type Subscription struct {
	EventName string
	Handler   func(ctx context.Context, event domain.DomainEvent) error
}

type EventBus interface {
	Publish(ctx context.Context, event domain.DomainEvent) error
	Subscribe(subscriber, eventName string, handler func(ctx context.Context, event domain.DomainEvent) error) (string, error)
	Unsubscribe(subscriptionId string) error
	Close()
}

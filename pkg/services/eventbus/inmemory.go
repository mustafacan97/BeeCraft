package event_bus

import (
	"context"
	"fmt"
	"platform/pkg/domain"

	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type inMemoryEventBus struct {
	subscriptions map[string]map[string]Subscription // eventName -> (SubscriptionId -> Subscription)
	mu            sync.RWMutex
}

func NewInMemoryEventBus() EventBus {
	return &inMemoryEventBus{
		subscriptions: make(map[string]map[string]Subscription),
	}
}

func (b *inMemoryEventBus) Publish(ctx context.Context, event domain.DomainEvent) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if subs, exists := b.subscriptions[event.GetEventName()]; exists {
		for _, sub := range subs {
			go sub.Handler(ctx, event)
		}
	}
	return nil
}

func (b *inMemoryEventBus) Subscribe(subscriber, eventName string, handler func(ctx context.Context, event domain.DomainEvent) error) (string, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	subId := uuid.New().String()
	if _, exists := b.subscriptions[eventName]; !exists {
		b.subscriptions[eventName] = make(map[string]Subscription)
	}
	b.subscriptions[eventName][subId] = Subscription{
		EventName: eventName,
		Handler:   handler,
	}

	return subId, nil
}

func (b *inMemoryEventBus) Unsubscribe(subscriptionId string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	for eventName, subs := range b.subscriptions {
		if _, exists := subs[subscriptionId]; exists {
			delete(subs, subscriptionId)
			if len(subs) == 0 {
				delete(b.subscriptions, eventName)
			}
			return nil
		}
	}
	return fmt.Errorf("unsubscribe failed: subscription with ID '%s' not found in any event", subscriptionId)
}

func (b *inMemoryEventBus) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	// You can clear memory by clearing all subscriptions.
	b.subscriptions = make(map[string]map[string]Subscription)
	zap.L().Info("InMemory event bus connection closed")
}

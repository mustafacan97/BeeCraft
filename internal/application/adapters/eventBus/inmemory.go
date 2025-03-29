package eventBus

import (
	"context"
	"errors"
	eventBus "platform/internal/application/ports/services"
	domainEvents "platform/internal/domain/events"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// InMemoryEventBus implements an event bus in-memory.
type inMemoryEventBus struct {
	subscriptions map[string]map[string]eventBus.Subscription // eventName -> (SubscriptionId -> Subscription)
	mu            sync.RWMutex
}

// NewInMemoryEventBus initialize a new InMemoryEventBus
func NewInMemoryEventBus() eventBus.EventBus {
	return &inMemoryEventBus{
		subscriptions: make(map[string]map[string]eventBus.Subscription),
	}
}

// Publish sends an event to all subscribers of the event name.
func (b *inMemoryEventBus) Publish(ctx context.Context, event domainEvents.Event) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if subs, exists := b.subscriptions[event.Name]; exists {
		for _, sub := range subs {
			go sub.Handler(ctx, event) // Execute asynchronously
		}
	}
	return nil
}

// Subscribe registers a new event handler.
func (b *inMemoryEventBus) Subscribe(subscriber, eventName string, handler func(ctx context.Context, event domainEvents.Event) error) (string, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	subId := uuid.New().String()
	if _, exists := b.subscriptions[eventName]; !exists {
		b.subscriptions[eventName] = make(map[string]eventBus.Subscription)
	}
	b.subscriptions[eventName][subId] = eventBus.Subscription{
		EventName: eventName,
		Handler:   handler,
	}

	return subId, nil
}

// Unsubscribe removes a subscriber by ID.
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
	return errors.New("subscription ID not found")
}

func (b *inMemoryEventBus) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	// You can clear memory by clearing all subscriptions.
	b.subscriptions = make(map[string]map[string]eventBus.Subscription)
	zap.L().Info("RabbitMQ connection closed")
}

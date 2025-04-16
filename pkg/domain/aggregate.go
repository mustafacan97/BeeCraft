package domain

import "github.com/google/uuid"

type AggregateRoot struct {
	BaseEntity[uuid.UUID]
	domainEvents []DomainEvent
}

func NewAggregateRoot(id uuid.UUID) AggregateRoot {
	return AggregateRoot{
		BaseEntity:   NewBaseEntityWithID(id),
		domainEvents: make([]DomainEvent, 0),
	}
}

func (a *AggregateRoot) AddEvent(event DomainEvent) {
	a.domainEvents = append(a.domainEvents, event)
}

// PullEvents → exports and cleans up events (for event dispatching)
func (a *AggregateRoot) PullEvents() []DomainEvent {
	events := a.domainEvents
	a.domainEvents = []DomainEvent{}
	return events
}

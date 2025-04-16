package domain

import "time"

type DomainEvent interface {
	GetEventName() string
	GetTimestamp() time.Time
	GetPayload() DomainEventPayload
}

type DomainEventPayload interface {
	Validate() error
}

type BaseDomainEvent struct {
	EventName string             `json:"event_name"`
	Timestamp time.Time          `json:"timestamp"`
	Payload   DomainEventPayload `json:"payload"`
}

func (e BaseDomainEvent) NewBaseDomainEvent(eventName string, payload DomainEventPayload) BaseDomainEvent {
	return BaseDomainEvent{
		EventName: eventName,
		Timestamp: time.Now(),
		Payload:   payload,
	}
}

func (e BaseDomainEvent) GetEventName() string {
	return e.EventName
}

func (e BaseDomainEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

func (e BaseDomainEvent) GetPayload() DomainEventPayload {
	return e.Payload
}

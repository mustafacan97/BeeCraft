package domain

import "time"

type Event interface {
	GetAggregateId() string
	GetAggregateType() string
	GetEventType() string
	GetTimestamp() time.Time
}

type BaseEvent struct {
	AggregateId   string    `json:"aggregate_id"`
	AggregateType string    `json:"aggregate_type"`
	EventType     string    `json:"event_type"`
	Timestamp     time.Time `json:"timestamp"`
	Payload       any       `json:"payload"`
}

func (e BaseEvent) GetAggregateId() string {
	return e.AggregateId
}

func (e BaseEvent) GetAggregateType() string {
	return e.AggregateType
}

func (e BaseEvent) GetEventType() string {
	return e.EventType
}

func (e BaseEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

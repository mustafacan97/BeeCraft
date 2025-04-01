package iam

import (
	"platform/internal/domain"
	"time"
)

type UserRegisteredEvent struct {
	Email string `json:"email"`
}

type UserPasswordChangedEvent struct {
	Email string `json:"email"`
}

func NewUserRegisteredEvent(aggregateId string, email string) domain.Event {
	return &domain.BaseEvent{
		AggregateId:   aggregateId,
		AggregateType: "USER",
		EventType:     "user.registered",
		Timestamp:     time.Now(),
		Payload: UserRegisteredEvent{
			Email: email,
		},
	}
}

func NewUserPasswordChangedEvent(aggregateId string, email string) domain.Event {
	return &domain.BaseEvent{
		AggregateId:   aggregateId,
		AggregateType: "USER",
		EventType:     "user.passwordChanged",
		Timestamp:     time.Now(),
		Payload: UserPasswordChangedEvent{
			Email: email,
		},
	}
}

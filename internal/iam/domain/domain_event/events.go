package domain_event

import (
	"platform/pkg/domain"
	"time"
)

type UserRegisteredEvent struct {
	Email string `json:"email"`
}

func (e UserRegisteredEvent) Validate() error {
	return nil
}

type UserPasswordChangedEvent struct {
	Email string `json:"email"`
}

func (e UserPasswordChangedEvent) Validate() error {
	return nil
}

func NewUserRegisteredEvent(aggregateId string, email string) domain.DomainEvent {
	return &domain.BaseDomainEvent{
		EventName: "user.registered",
		Timestamp: time.Now(),
		Payload: UserRegisteredEvent{
			Email: email,
		},
	}
}

func NewUserPasswordChangedEvent(aggregateId string, email string) domain.DomainEvent {
	return &domain.BaseDomainEvent{
		EventName: "user.passwordChanged",
		Timestamp: time.Now(),
		Payload: UserPasswordChangedEvent{
			Email: email,
		},
	}
}

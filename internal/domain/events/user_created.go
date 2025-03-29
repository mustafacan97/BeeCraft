package domainEvents

import "time"

// UserCreatedDomainEvent embeds BaseDomainEvent
type userCreatedDomainEvent struct {
	Email string `json:"email"`
}

func NewUserCreatedDomainEvent(email string) *Event {
	return &Event{
		Name: "user.created",
		Payload: userCreatedDomainEvent{
			Email: email,
		},
		Timestamp: time.Now(),
	}
}

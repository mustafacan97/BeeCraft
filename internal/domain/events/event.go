package domainEvents

import "time"

// Event represents a message passed through the EventBus.
type Event struct {
	Name      string
	Payload   any
	Timestamp time.Time
}

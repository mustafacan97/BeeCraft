package event_notification

import (
	"time"

	"github.com/google/uuid"
)

type EmailAccountUpdatedEvent struct {
	ProjectID uuid.UUID `json:"project_id"`
	Email     string    `json:"email"`
	UpdatedAt time.Time `json:"deleted_at"`
}

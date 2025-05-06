package event_notification

import (
	"time"

	"github.com/google/uuid"
)

type EmailAccountDeletedEvent struct {
	ProjectID uuid.UUID `json:"project_id"`
	Email     string    `json:"email"`
	DeletedAt time.Time `json:"deleted_at"`
}

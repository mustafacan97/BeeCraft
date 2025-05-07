package event_notification

import "github.com/google/uuid"

type EmailAccountDeletedEvent struct {
	ProjectID uuid.UUID `json:"project_id"`
	Email     string    `json:"email"`
}

func NewEmailAccountDeletedEvent(projectID uuid.UUID, email string) EmailAccountDeletedEvent {
	return EmailAccountDeletedEvent{
		ProjectID: projectID,
		Email:     email,
	}
}

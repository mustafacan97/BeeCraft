package event_notification

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type EmailAccountCreatedEvent struct {
	ProjectID uuid.UUID `json:"project_id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func NewEmailAccountCreatedEvent(projectID uuid.UUID, email string) EmailAccountCreatedEvent {
	return EmailAccountCreatedEvent{
		ProjectID: projectID,
		Email:     email,
		CreatedAt: time.Now(),
	}
}

type EmailAccountCreatedEventHandler struct{}

func (c *EmailAccountCreatedEventHandler) Handle(ctx context.Context, event *EmailAccountCreatedEvent) error {
	//Do something with the event here !
	zap.L().Info(fmt.Sprintf("Test message: %s", event.Email))
	return nil
}

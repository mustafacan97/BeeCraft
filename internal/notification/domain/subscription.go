package domain

import (
	"platform/internal/shared"
	"platform/pkg/domain/valueobject"
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	projectId uuid.UUID
	email     valueobject.Email
	language  shared.Language
	createdAt time.Time
}

func NewSubscription(projectId uuid.UUID, email valueobject.Email, language shared.Language) *Subscription {
	return &Subscription{
		projectId: projectId,
		email:     email,
		language:  language,
		createdAt: time.Now(),
	}
}

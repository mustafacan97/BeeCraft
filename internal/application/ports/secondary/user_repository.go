package secondary

import (
	"context"

	"platform/internal/domain"

	"github.com/google/uuid"
)

type UserRepository interface {
	// Read operations
	GetById(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)

	// Write operations
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

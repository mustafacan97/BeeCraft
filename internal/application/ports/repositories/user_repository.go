package repositories

import (
	"context"

	"platform/internal/domain/iam"

	"github.com/google/uuid"
)

type UserRepository interface {
	// Read operations
	GetById(ctx context.Context, id uuid.UUID) (*iam.User, error)
	GetByEmail(ctx context.Context, email string) (*iam.User, error)
	Exists(ctx context.Context, email string) (bool, error)

	// Write operations
	Create(ctx context.Context, user *iam.User) error
	Update(ctx context.Context, user *iam.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

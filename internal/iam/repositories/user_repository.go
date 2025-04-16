package repositories

import (
	"context"
	"platform/internal/iam/domain"

	"github.com/google/uuid"
)

type UserRepository interface {
	// QUERY
	GetById(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Exists(ctx context.Context, email string) (bool, error)

	// COMMAND
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

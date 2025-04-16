package repositories

import (
	"context"
	"platform/internal/notification/domain"

	"github.com/google/uuid"
)

type EmailAccountRepository interface {
	// QUERY
	GetAll(ctx context.Context) ([]*domain.EmailAccount, error)
	GetByEmail(ctx context.Context, email string) (*domain.EmailAccount, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.EmailAccount, error)

	// COMMAND
	Create(ctx context.Context, account *domain.EmailAccount) error
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, account *domain.EmailAccount) error
}

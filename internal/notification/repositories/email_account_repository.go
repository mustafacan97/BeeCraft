package repositories

import (
	"context"
	"platform/internal/notification/domain"
	"platform/pkg/domain/valueobject"

	"github.com/google/uuid"
)

type EmailAccountRepository interface {
	// QUERY
	GetAll(ctx context.Context, page, pageSize int) ([]*domain.EmailAccount, int, error)
	GetByEmail(ctx context.Context, email valueobject.Email) (*domain.EmailAccount, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.EmailAccount, error)

	// COMMAND
	Create(ctx context.Context, account *domain.EmailAccount) error
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, account *domain.EmailAccount) error
}

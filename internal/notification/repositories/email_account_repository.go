package repositories

import (
	"context"
	"platform/internal/notification/domain"
	vo "platform/pkg/domain/value_object"
)

type EmailAccountRepository interface {
	// QUERY
	GetAll(ctx context.Context) ([]*domain.EmailAccount, error)
	GetByEmail(ctx context.Context, email vo.Email) (*domain.EmailAccount, error)

	// COMMAND
	Create(ctx context.Context, account *domain.EmailAccount) error
	Delete(ctx context.Context, email vo.Email) error
	Update(ctx context.Context, account *domain.EmailAccount) error
}

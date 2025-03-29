package repositories

import (
	"context"

	"platform/internal/domain"
)

type RoleRepository interface {
	// Read operations
	GetById(ctx context.Context, id int) (*domain.Role, error)
	GetByProjectId(ctx context.Context, projectId string) ([]*domain.Role, error)
	GetSystemRoleByName(ctx context.Context, name string) (*domain.Role, error)

	// Write operations
	Create(ctx context.Context, role *domain.Role) error
	Update(ctx context.Context, role *domain.Role) error
	Delete(ctx context.Context, id int) error
}

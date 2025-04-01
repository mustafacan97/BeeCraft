package repositories

import (
	"context"

	"platform/internal/domain/iam"
)

type RoleRepository interface {
	// Read operations
	GetById(ctx context.Context, id int) (*iam.Role, error)
	GetByProjectId(ctx context.Context, projectId string) ([]*iam.Role, error)
	GetSystemRoleByName(ctx context.Context, name string) (*iam.Role, error)

	// Write operations
	Create(ctx context.Context, role *iam.Role) error
	Update(ctx context.Context, role *iam.Role) error
	Delete(ctx context.Context, id int) error
}

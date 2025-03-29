package postgresql

import (
	"context"
	"fmt"
	"platform/internal/application/ports/repositories"
	"platform/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type roleRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewRoleRepository(pool *pgxpool.Pool) repositories.RoleRepository {
	return &roleRepositoryImpl{pool: pool}
}

func (r *roleRepositoryImpl) Create(ctx context.Context, role *domain.Role) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO roles (name, project_id) VALUES ($1, $2)", role.Name, role.ProjectId)
	return err
}

func (r *roleRepositoryImpl) GetById(ctx context.Context, id int) (*domain.Role, error) {
	var role domain.Role
	err := r.pool.QueryRow(ctx, "SELECT id, name, project_id FROM roles WHERE id = $1", id).Scan(&role.Id, &role.Name, &role.ProjectId)
	return &role, err
}

func (r *roleRepositoryImpl) GetByProjectId(ctx context.Context, projectId string) ([]*domain.Role, error) {
	var roles []*domain.Role
	rows, err := r.pool.Query(ctx, "SELECT id, name, project_id FROM roles WHERE project_id = $1", projectId)

	if err != nil {
		return roles, err
	}

	for rows.Next() {
		var role domain.Role
		if err := rows.Scan(&role.Id, &role.Name, &role.ProjectId); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		roles = append(roles, &role)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return roles, nil
}

func (r *roleRepositoryImpl) GetSystemRoleByName(ctx context.Context, name string) (*domain.Role, error) {
	var role domain.Role
	sql := "SELECT id, name, project_id FROM roles WHERE project_id IS NULL AND name = $1"
	err := r.pool.QueryRow(ctx, sql, name).Scan(&role.Id, &role.Name, nil)

	if err != nil {
		// The specified role not found
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error scanning role: %w", err)
	}

	return &role, nil
}

func (r *roleRepositoryImpl) Update(ctx context.Context, role *domain.Role) error {
	_, err := r.pool.Exec(ctx, `UPDATE roles SET name = $1`, &role.Name)
	return err
}

func (r *roleRepositoryImpl) Delete(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM roles WHERE id = $1`, id)
	return err
}

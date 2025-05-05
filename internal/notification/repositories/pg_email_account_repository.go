package repositories

import (
	"context"
	"fmt"
	"platform/internal/notification/domain"
	"platform/internal/shared"
	"platform/pkg/domain/valueobject"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrMissingProjectID = fmt.Errorf("project ID not found in context")
	ErrInvalidProjectID = fmt.Errorf("project ID has invalid type")
	ErrEmptyProjectID   = fmt.Errorf("project ID is empty (uuid.Nil)")
)

type pgEmailAccountRepository struct {
	pool *pgxpool.Pool
}

func NewPgEmailAccountRepository(pool *pgxpool.Pool) EmailAccountRepository {
	return &pgEmailAccountRepository{pool: pool}
}

// QUERY
func (p *pgEmailAccountRepository) GetAll(ctx context.Context, page, pageSize int) ([]*domain.EmailAccount, int, error) {
	projectID, _ := ctx.Value(shared.ProjectIDContextKey).(uuid.UUID)
	offset := (page - 1) * pageSize

	var totalCount int
	countSQL := "SELECT COUNT(*) FROM notification.email_accounts WHERE project_id = $1"
	if err := p.pool.QueryRow(ctx, countSQL, projectID).Scan(&totalCount); err != nil {
		return nil, 0, err
	}

	dataSQL := `SELECT * FROM notification.email_accounts WHERE project_id = $1 ORDER BY created_at LIMIT $2 OFFSET $3`
	rows, err := p.pool.Query(ctx, dataSQL, projectID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[EmailAccountDTO])
	if err != nil {
		return nil, 0, err
	}

	accounts := make([]*domain.EmailAccount, 0, len(dtos))
	for _, dto := range dtos {
		accounts = append(accounts, dto.ToDomain())
	}

	return accounts, totalCount, nil
}

func (p *pgEmailAccountRepository) GetByEmail(ctx context.Context, email valueobject.Email) (*domain.EmailAccount, error) {
	projectID, _ := ctx.Value(shared.ProjectIDContextKey).(uuid.UUID)
	sql := "SELECT * FROM notification.email_accounts WHERE project_id = $1 AND email = $2"
	rows, err := p.pool.Query(ctx, sql, projectID, email.GetValue())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dto, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[EmailAccountDTO])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return dto.ToDomain(), nil
}

func (p *pgEmailAccountRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.EmailAccount, error) {
	projectID, _ := ctx.Value(shared.ProjectIDContextKey).(uuid.UUID)
	sql := "SELECT * FROM notification.email_accounts WHERE project_id = $1 AND id = $2"
	rows, err := p.pool.Query(ctx, sql, projectID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dto, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[EmailAccountDTO])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return dto.ToDomain(), nil
}

// COMMAND
func (p *pgEmailAccountRepository) Create(ctx context.Context, account *domain.EmailAccount) error {
	query := `
		INSERT INTO notification.email_accounts (
			id,
			project_id, 
			email, 
			display_name,
			host, 
			port, 
			enable_ssl,
			type_id,
			username,
			password,
			client_id,
			client_secret,
			tenant_id,
			access_token,
			refresh_token,
			expire_at,
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`

	_, err := p.pool.Exec(ctx, query, ToDTO(account).ToValues()...)
	return err
}

func (p *pgEmailAccountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	projectID, _ := ctx.Value(shared.ProjectIDContextKey).(uuid.UUID)
	sql := "DELETE FROM notification.email_accounts WHERE project_id = $1 AND id = $2"
	_, err := p.pool.Exec(ctx, sql, projectID, id)
	return err
}

func (p *pgEmailAccountRepository) Update(ctx context.Context, account *domain.EmailAccount) error {
	query := `
		UPDATE notification.email_accounts SET
			email = $3,
			display_name = $4,
			host = $5,
			port = $6,
			enable_ssl = $7,
			type_id = $8,			
			username = $9,
			password = $10,
			client_id = $11,
			tenant_id = $12,
			client_secret = $13,
			access_token = $14,
			refresh_token = $15,
			expire_at = $16
		WHERE id = $1 AND project_id = $2
	`

	_, err := p.pool.Exec(ctx, query, ToDTO(account).ToValues()[0:16]...)
	return err
}

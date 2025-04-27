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
func (p *pgEmailAccountRepository) GetAll(ctx context.Context) ([]*domain.EmailAccount, error) {
	projectID, err := getProjectID(ctx)
	if err != nil {
		return []*domain.EmailAccount{}, err
	}

	sql := "SELECT * FROM email_accounts WHERE project_id = $1"
	rows, err := p.pool.Query(ctx, sql, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[EmailAccountDTO])
	if err != nil {
		return []*domain.EmailAccount{}, err
	}

	accounts := make([]*domain.EmailAccount, 0, len(dtos))
	for _, dto := range dtos {
		accounts = append(accounts, dto.ToDomain())
	}

	return accounts, nil
}

func (p *pgEmailAccountRepository) GetByEmail(ctx context.Context, email valueobject.Email) (*domain.EmailAccount, error) {
	projectID, _ := ctx.Value(shared.ProjectIDContextKey).(string)
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
	projectID, err := getProjectID(ctx)
	if err != nil {
		return nil, err
	}

	sql := "SELECT * FROM email_accounts WHERE project_id = $1 AND id = $2"
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
	projectID, err := getProjectID(ctx)
	if err != nil {
		return err
	}

	query := "DELETE FROM email_accounts WHERE id = $1 AND project_id = $2"

	_, err = p.pool.Exec(ctx, query, id, projectID)
	return err
}

func (p *pgEmailAccountRepository) Update(ctx context.Context, account *domain.EmailAccount) error {
	dto := ToDTO(account)

	query := `
		UPDATE email_accounts SET
			email = $1,
			display_name = $2,
			host = $3,
			port = $4,
			type_id = $5,
			enable_ssl = $6,
			username = $7,
			password = $8,
			client_id = $9,
			client_secret = $10,
			tenant_id = $11,
			access_token = $12,
			refresh_token = $13,
			expire_at = $14
		WHERE id = $15 AND project_id = $16
	`

	values := []any{
		dto.Email,
		dto.DisplayName,
		dto.Host,
		dto.Port,
		dto.TypeID,
		dto.EnableSsl,
		dto.Username,
		dto.Password,
		dto.ClientID,
		dto.ClientSecret,
		dto.TenantID,
		dto.AccessToken,
		dto.RefreshToken,
		dto.ExpireAt,
		dto.ID,
		dto.ProjectID,
	}

	_, err := p.pool.Exec(ctx, query, values...)
	return err
}

func getProjectID(ctx context.Context) (uuid.UUID, error) {
	val := ctx.Value(shared.ProjectIDContextKey)
	if val == nil {
		return uuid.Nil, ErrMissingProjectID
	}

	projectID, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil, ErrInvalidProjectID
	}

	if projectID == uuid.Nil {
		return uuid.Nil, ErrEmptyProjectID
	}

	return projectID, nil
}

package repositories

import (
	"context"
	"fmt"
	"platform/internal/notification/domain"

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
		account, err := dto.ToDomain()
		if err != nil {
			// TODO: log the error
		} else {
			accounts = append(accounts, account)
		}
	}

	return accounts, nil
}

func (p *pgEmailAccountRepository) GetByEmail(ctx context.Context, email string) (*domain.EmailAccount, error) {
	projectID, err := getProjectID(ctx)
	if err != nil {
		return nil, err
	}

	sql := "SELECT * FROM email_accounts WHERE project_id = $1 AND email = $2"
	rows, err := p.pool.Query(ctx, sql, projectID, email)
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

	emailAccount, err := dto.ToDomain()
	if err != nil {
		return nil, err
	}

	return emailAccount, nil

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

	emailAccount, err := dto.ToDomain()
	if err != nil {
		return nil, err
	}

	return emailAccount, nil
}

// COMMAND
func (p *pgEmailAccountRepository) Create(ctx context.Context, account *domain.EmailAccount) error {
	dto := ToDTO(account)

	query := `
		INSERT INTO email_accounts (
			id, project_id, email, display_name, host, port, type_id, enable_ssl, created_at,
			username, password, 
			client_id, client_secret, tenant_id,
			access_token, refresh_token, expire_at
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
			$10, $11, 
			$12, $13, $14, 
			$15, $16, $17
		)
	`

	_, err := p.pool.Exec(ctx, query, dto.ToValues()...)
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
		dto.email,
		dto.displayName,
		dto.host,
		dto.port,
		dto.typeID,
		dto.enableSsl,
		dto.username,
		dto.password,
		dto.clientID,
		dto.clientSecret,
		dto.tenantID,
		dto.accessToken,
		dto.refreshToken,
		dto.expireAt,
		dto.id,
		dto.projectID,
	}

	_, err := p.pool.Exec(ctx, query, values...)
	return err
}

func getProjectID(ctx context.Context) (uuid.UUID, error) {
	val := ctx.Value("projectID")
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

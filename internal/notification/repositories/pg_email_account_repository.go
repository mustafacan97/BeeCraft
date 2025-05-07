package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"platform/internal/notification/domain"
	"platform/internal/shared"
	vo "platform/pkg/domain/value_object"
	"platform/pkg/services/cache"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	ErrMissingProjectID = fmt.Errorf("project ID not found in context")
	ErrInvalidProjectID = fmt.Errorf("project ID has invalid type")
	ErrEmptyProjectID   = fmt.Errorf("project ID is empty (uuid.Nil)")
)

type pgEmailAccountRepository struct {
	pool  *pgxpool.Pool
	cache cache.CacheManager
}

func NewPgEmailAccountRepository(pool *pgxpool.Pool, cache cache.CacheManager) EmailAccountRepository {
	return &pgEmailAccountRepository{
		pool:  pool,
		cache: cache,
	}
}

// QUERY
func (p *pgEmailAccountRepository) GetAll(ctx context.Context) ([]*domain.EmailAccount, error) {
	// STEP-1: Get project identifier and validate
	pidVal := ctx.Value(shared.ProjectIDContextKey)
	projectID, ok := pidVal.(uuid.UUID)
	if !ok {
		return nil, shared.ErrInvalidContext
	}

	// STEP-2: Create a cache key
	cacheKey := cache.CacheKey{
		Key:  cacheKeyAll(projectID),
		Time: cache.DefaultTTL,
	}

	// STEP-4: Check if the data is in the cache service
	cached, err := p.cache.Get(ctx, cacheKey)
	if err != nil {
		zap.L().Warn("cache GET error", zap.Error(err), zap.String("key", cacheKey.Key))
	} else if cached != "" {
		var dtoList []EmailAccountDTO
		if err := json.Unmarshal([]byte(cached), &dtoList); err == nil {
			accounts := make([]*domain.EmailAccount, 0, len(dtoList))
			for _, dto := range dtoList {
				accounts = append(accounts, dto.ToDomain())
			}
			return accounts, nil
		}

		// Clear any corrupted data
		p.clearCaches(ctx, cacheKey.Key)
		zap.L().Warn("cache unmarshal failed, key removed", zap.String("key", cacheKey.Key), zap.Error(err))
	}

	// STEP-5: Get result from database
	sql := `SELECT * FROM notification.email_accounts WHERE project_id = $1 ORDER BY created_at`
	rows, err := p.pool.Query(ctx, sql, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dtoList, err := pgx.CollectRows(rows, pgx.RowToStructByName[EmailAccountDTO])
	if err != nil {
		return nil, err
	}

	// STEP-6: Convert from dto to domain
	accounts := make([]*domain.EmailAccount, 0, len(dtoList))
	for _, dto := range dtoList {
		accounts = append(accounts, dto.ToDomain())
	}

	// STEP-7: Save result the cache service
	if serialized, err := json.Marshal(dtoList); err == nil {
		if err := p.cache.Set(ctx, cacheKey, string(serialized)); err != nil {
			zap.L().Error("an error occurred while writing to cache", zap.Error(err))
		}
	}

	return accounts, nil
}

func (p *pgEmailAccountRepository) GetByEmail(ctx context.Context, email vo.Email) (*domain.EmailAccount, error) {
	// STEP-1: Get project identifier and validate
	pidVal := ctx.Value(shared.ProjectIDContextKey)
	projectID, ok := pidVal.(uuid.UUID)
	if !ok {
		return nil, shared.ErrInvalidContext
	}

	// STEP-2: Create a cache key
	cacheKey := cache.CacheKey{
		Key:  cacheKeyByEmail(projectID, email),
		Time: cache.DefaultTTL,
	}

	// STEP-3: Create default DTO object
	var dto EmailAccountDTO

	// STEP-4: Check if the data is in the cache service
	cached, err := p.cache.Get(ctx, cacheKey)
	if err != nil {
		zap.L().Warn("cache GET error", zap.Error(err), zap.String("key", cacheKey.Key))
	} else if cached != "" {
		if err := json.Unmarshal([]byte(cached), &dto); err == nil {
			return dto.ToDomain(), nil
		}

		// Clear any corrupted data
		p.clearCaches(ctx, cacheKey.Key)
		zap.L().Warn("cache unmarshal failed, key removed", zap.String("key", cacheKey.Key), zap.Error(err))
	}

	// STEP-5: Get result from database
	sql := "SELECT * FROM notification.email_accounts WHERE project_id = $1 AND email = $2"
	rows, err := p.pool.Query(ctx, sql, projectID, email.Value())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dto, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[EmailAccountDTO])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// STEP-6: Save result the cache service
	if serialized, err := json.Marshal(dto); err == nil {
		err = p.cache.Set(ctx, cacheKey, string(serialized))
		if err != nil {
			zap.L().Error("an error occurred while writing to cache", zap.Error(err))
		}
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
func (p *pgEmailAccountRepository) Create(ctx context.Context, ea *domain.EmailAccount) error {
	// STEP-1: Get project identifier and validate
	pidVal := ctx.Value(shared.ProjectIDContextKey)
	projectID, ok := pidVal.(uuid.UUID)
	if !ok {
		zap.L().Error("project ID not found in context", zap.Any("value", pidVal))
		return shared.ErrInvalidContext
	}
	ea.SetProjectID(projectID)

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

	dto := EmailAccountDTO{}
	_, err := p.pool.Exec(ctx, query, dto.ToDTO(ea).GetValues()...)
	return err
}

func (p *pgEmailAccountRepository) Delete(ctx context.Context, email vo.Email) error {
	// STEP-1: Get project identifier and validate
	pidVal := ctx.Value(shared.ProjectIDContextKey)
	projectID, ok := pidVal.(uuid.UUID)
	if !ok {
		return shared.ErrInvalidContext
	}

	// STEP-2: Delete from database
	sql := "DELETE FROM notification.email_accounts WHERE project_id = $1 AND email = $2"
	_, err := p.pool.Exec(ctx, sql, projectID, email.Value())
	if err != nil {
		return fmt.Errorf("failed to delete email account: %w", err)
	}

	// STEP-3: Remove related caches
	p.clearCaches(ctx, cacheKeyByEmail(projectID, email), cacheKeyAll(projectID))
	err = p.cache.Remove(ctx, cacheKeyByEmail(projectID, email))
	if err != nil {
		zap.L().Warn("an error occurred while removing cache key", zap.Error(err))
	}
	return nil
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
	dto := EmailAccountDTO{}
	_, err := p.pool.Exec(ctx, query, dto.ToDTO(account).GetValues()[0:16]...)
	return err
}

func (p *pgEmailAccountRepository) clearCaches(ctx context.Context, keys ...string) {
	for _, key := range keys {
		err := p.cache.Remove(ctx, key)
		if err != nil {
			zap.L().Warn("an error occurred while removing cache key", zap.Error(err))
		}
	}
}

func cacheKeyByEmail(projectID uuid.UUID, email vo.Email) string {
	return fmt.Sprintf("notification:email_accounts:%s:%s", projectID.String(), email.Value())
}

func cacheKeyAll(projectID uuid.UUID) string {
	return fmt.Sprintf("notification:email_accounts:%s", projectID.String())
}

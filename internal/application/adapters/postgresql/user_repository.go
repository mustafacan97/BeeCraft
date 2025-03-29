package postgresql

import (
	"context"
	"errors"
	"fmt"
	"platform/internal/application/ports/repositories"
	"platform/internal/domain"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) repositories.UserRepository {
	return &userRepositoryImpl{pool: pool}
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *domain.User) error {
	if len(user.Roles) == 0 {
		return errors.New("user must have at least a role")
	}

	// Construct parameterized SQL for bulk insert
	valueStrings := make([]string, len(user.Roles))
	valueArgs := make([]any, 0, len(user.Roles)*2)

	for i, role := range user.Roles {
		paramIdx := i * 2
		valueStrings[i] = fmt.Sprintf("($%d, $%d)", paramIdx+1, paramIdx+2)
		valueArgs = append(valueArgs, user.Id.String(), role.Id)
	}

	userSql := "INSERT INTO users (id, email, email_validated, phone_validated, password_hash, failed_login_attempts, is_system_user, created_at, active, deleted) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
	roleSql := "INSERT INTO user_role_mappings (user_id, role_id) VALUES " + strings.Join(valueStrings, ",")

	return runInTransaction(ctx, r.pool, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, userSql, user.Id, user.Email, false, false, user.PasswordHash, 0, false, time.Now(), true, false)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, roleSql, valueArgs...)
		return err
	})
}

func (r *userRepositoryImpl) GetById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	sql := `SELECT * FROM users WHERE id = $1`
	err := r.pool.QueryRow(ctx, sql, id).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.EmailValidated,
		&user.Phone,
		&user.PhoneValidated,
		&user.Gender,
		&user.BirthDate,
		&user.PasswordHash,
		&user.LastPasswordChangeAt,
		&user.FailedLoginAttempts,
		&user.CannotLoginUntilAt,
		&user.RefreshToken,
		&user.RefreshTokenExpireAt,
		&user.LastIpAddress,
		&user.LastLoginAt,
		&user.IsSystemUser,
		&user.AdminComment,
		&user.CreatedAt,
		&user.Active,
		&user.Deleted,
	)
	return &user, err
}

func (r *userRepositoryImpl) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	sql := `SELECT * FROM users WHERE email = $1`
	err := r.pool.QueryRow(ctx, sql, email).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.EmailValidated,
		&user.Phone,
		&user.PhoneValidated,
		&user.Gender,
		&user.BirthDate,
		&user.PasswordHash,
		&user.LastPasswordChangeAt,
		&user.FailedLoginAttempts,
		&user.CannotLoginUntilAt,
		&user.RefreshToken,
		&user.RefreshTokenExpireAt,
		&user.LastIpAddress,
		&user.LastLoginAt,
		&user.IsSystemUser,
		&user.AdminComment,
		&user.CreatedAt,
		&user.Active,
		&user.Deleted,
	)
	return &user, err
}

func (r *userRepositoryImpl) Update(ctx context.Context, user *domain.User) error {
	sql := `
		UPDATE users
		SET 
			first_name = $1, 
			last_name = $2, 
			email = $3,
			email_validate = $4,
			phone = $5,
			phone_validated = $6,
			gender = $7,
			birth_date = $8,
			password_hash = $9,
			last_password_change_at = $10,
			failed_login_attempts = $11,
			cannot_login_until_at = $12,
			refresh_token = $13,
			refresh_token_expire_at = $14,
			last_ip_address = $15,
			last_login_at = $16,
			is_system_user = $17,
			admin_comment = $18,
			active = $19,
			deleted = $20
	`
	_, err := r.pool.Exec(
		ctx,
		sql,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.EmailValidated,
		&user.Phone,
		&user.PhoneValidated,
		&user.Gender,
		&user.BirthDate,
		&user.PasswordHash,
		&user.LastPasswordChangeAt,
		&user.FailedLoginAttempts,
		&user.CannotLoginUntilAt,
		&user.RefreshToken,
		&user.RefreshTokenExpireAt,
		&user.LastIpAddress,
		&user.LastLoginAt,
		&user.IsSystemUser,
		&user.AdminComment,
		&user.Active,
		&user.Deleted)

	return err
}

func (r *userRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `DELETE FROM users WHERE id = $1`
	_, err := r.pool.Exec(ctx, sql, id)
	return err
}

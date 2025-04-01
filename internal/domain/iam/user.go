package iam

import (
	"errors"
	"time"

	"platform/internal/domain"
	"platform/internal/enum"

	"github.com/google/uuid"
)

var (
	ErrInvalidPassword    = errors.New("invalid password")
	ErrInvalidUser        = errors.New("invalid user")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid username or password")
)

const (
	max_failed_attempts = 5
)

type User struct {
	domain.BaseAggregateRoot
	FirstName            *string
	LastName             *string
	Email                string
	EmailValidated       bool
	Phone                *string
	PhoneValidated       bool
	Gender               enum.GenderType
	BirthDate            *time.Time
	PasswordHash         string
	LastPasswordChangeAt *time.Time
	FailedLoginAttempts  uint8
	CannotLoginUntilAt   *time.Time
	RefreshToken         *string
	RefreshTokenExpireAt *time.Time
	LastIpAddress        *string
	LastLoginAt          *time.Time
	IsSystemUser         bool
	AdminComment         *string
	CreatedAt            time.Time
	UpdatedAt            *time.Time
	Active               bool
	Deleted              bool
	Roles                []Role
}

func NewUser(email, rawPassword string, roles []Role) (*User, error) {
	password, err := NewPassword(rawPassword)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := password.Hash()
	if err != nil {
		return nil, err
	}

	user := &User{
		BaseAggregateRoot: domain.BaseAggregateRoot{
			Id: uuid.New(),
		},
		FirstName:            nil,
		LastName:             nil,
		Email:                email,
		EmailValidated:       false,
		Phone:                nil,
		PhoneValidated:       false,
		Gender:               enum.NotSpecified,
		BirthDate:            nil,
		PasswordHash:         hashedPassword,
		LastPasswordChangeAt: nil,
		FailedLoginAttempts:  0,
		CannotLoginUntilAt:   nil,
		RefreshToken:         nil,
		RefreshTokenExpireAt: nil,
		LastIpAddress:        nil,
		LastLoginAt:          nil,
		IsSystemUser:         false,
		AdminComment:         nil,
		CreatedAt:            time.Now(),
		UpdatedAt:            nil,
		Active:               true,
		Deleted:              false,
		Roles:                roles,
	}

	return user, nil
}

func (u *User) Authenticate(rawPassword string) bool {
	password := Password(rawPassword)
	return password.Matches(u.PasswordHash)
}

func (u *User) ChangePassword(rawOldPassword, rawNewPassword string) error {
	ok := u.Authenticate(rawOldPassword)
	if !ok {
		return ErrInvalidPassword
	}

	newPassword, err := NewPassword(rawNewPassword)
	if err != nil {
		return err
	}

	_, err = newPassword.Hash()
	if err != nil {
		return err
	}

	return nil
}

func (u *User) UpdateEmail(newEmail string) {
	u.Email = newEmail
	u.EmailValidated = false
}

func (u *User) UpdatePhone(newPhone string) {
	u.Phone = &newPhone
	u.PhoneValidated = false
}

func (u *User) IncrementFailedLoginAttempts() {
	u.FailedLoginAttempts++
	if u.FailedLoginAttempts >= max_failed_attempts {
		lockTime := time.Now().Add(30 * time.Minute) // Lock the user out for 30 minutes after 5 failed attempts
		u.CannotLoginUntilAt = &lockTime
	}
}

func (u *User) ResetFailedLoginAttempts() {
	u.FailedLoginAttempts = 0
	u.CannotLoginUntilAt = nil
}

func (u *User) SetLastLogin(ipAddress string) {
	now := time.Now()
	u.LastLoginAt = &now
	u.LastIpAddress = &ipAddress
}

func (u *User) MarkAsDeleted() {
	u.Deleted = true
	u.Active = false
}

func (u *User) ActivateUser() {
	now := time.Now()
	u.Active = true
	u.Deleted = false
	u.LastLoginAt = &now
}

package domain

import (
	"time"

	"platform/internal/enum"
	services "platform/internal/service"

	"github.com/google/uuid"
)

const (
	max_failed_attempts = 5
)

type User struct {
	Id                   uuid.UUID
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
	Active               bool
	Deleted              bool
	Roles                []Role
}

func NewUser(email, password string, roles []Role) (*User, error) {
	passwordHasher := &services.PasswordHasher{}

	// Hash and salt the password
	hashedPassword, err := passwordHasher.HashWithBcrypt(password)
	if err != nil {
		return nil, err // Return error if password hashing fails
	}

	return &User{
		Id:                   uuid.New(),
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
		Active:               true,
		Deleted:              false,
		Roles:                roles,
	}, nil
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

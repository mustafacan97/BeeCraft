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
	FirstName            string
	LastName             string
	Email                string
	EmailValidated       bool
	Phone                string
	PhoneValidated       bool
	Gender               enum.GenderType
	BirthDate            time.Time
	PasswordHash         string
	LastPasswordChangeAt time.Time
	FailedLoginAttempts  uint8
	CannotLoginUntilAt   time.Time
	RefreshToken         string
	RefreshTokenExpireAt time.Time
	LastIpAddress        string
	LastLoginAt          time.Time
	IsSystemUser         bool
	AdminComment         string
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
		Id:                   uuid.New(),        // Generate a new UUID for the user
		FirstName:            "",                // Empty by default, can be updated later
		LastName:             "",                // Empty by default, can be updated later
		Email:                email,             // Provided email
		EmailValidated:       false,             // Email validation flag, default to false
		Phone:                "",                // Empty by default, can be updated later
		PhoneValidated:       false,             // Phone validation flag, default to false
		Gender:               enum.NotSpecified, // Default value
		BirthDate:            time.Time{},       // No birth date set by default
		PasswordHash:         hashedPassword,    // Provided password hash
		LastPasswordChangeAt: time.Now(),        // Set to the current time
		FailedLoginAttempts:  0,                 // Default to 0 failed login attempts
		CannotLoginUntilAt:   time.Time{},       // No restriction by default
		RefreshToken:         "",                // Empty by default
		RefreshTokenExpireAt: time.Time{},       // No expiration set
		LastIpAddress:        "",
		LastLoginAt:          time.Time{},
		IsSystemUser:         false,
		AdminComment:         "",
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
	u.Phone = newPhone
	u.PhoneValidated = false
}

func (u *User) IncrementFailedLoginAttempts() {
	u.FailedLoginAttempts++
	if u.FailedLoginAttempts >= max_failed_attempts {
		u.CannotLoginUntilAt = time.Now().Add(30 * time.Minute) // Lock the user out for 30 minutes after 5 failed attempts
	}
}

func (u *User) ResetFailedLoginAttempts() {
	u.FailedLoginAttempts = 0
	u.CannotLoginUntilAt = time.Time{}
}

func (u *User) SetLastLogin(ipAddress string) {
	u.LastLoginAt = time.Now()
	u.LastIpAddress = ipAddress
}

func (u *User) MarkAsDeleted() {
	u.Deleted = true
	u.Active = false
}

func (u *User) ActivateUser() {
	u.Active = true
	u.Deleted = false
	u.LastLoginAt = time.Now()
}

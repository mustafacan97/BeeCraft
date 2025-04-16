package valueobject

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordLength    = errors.New("password must be between 8 and 16 characters long")
	ErrPasswordLowerCase = errors.New("password must have at least one lowercase letter")
	ErrPasswordUpperCase = errors.New("password must have at least one uppercase letter")
	ErrPasswordDigit     = errors.New("password must contain at least one number")
	ErrHashingPassword   = errors.New("hashing failed")
)

type Password string

func NewPassword(rawPassword string) (Password, error) {
	if len(rawPassword) < 8 || len(rawPassword) > 16 {
		return "", ErrPasswordLength
	}

	hasLower := regexp.MustCompile(`[a-z]`).MatchString(rawPassword)
	if !hasLower {
		return "", ErrPasswordLowerCase
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(rawPassword)
	if !hasUpper {
		return "", ErrPasswordUpperCase
	}

	hasDigit := regexp.MustCompile(`\d`).MatchString(rawPassword)
	if !hasDigit {
		return "", ErrPasswordDigit
	}

	return Password(rawPassword), nil
}

// HashWithBcrypt hashes the password using bcrypt (recommended for security).
// Built-in salting mechanism
func (p Password) Hash() (string, error) {
	// Generate bcrypt hash of the password with a default cost factor (10)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return "", ErrHashingPassword
	}
	return string(hashedPassword), nil
}

// VerifyBcryptPassword verifies if the given raw password matches the bcrypt hash.
func (p Password) Matches(storedHash string) bool {
	// Compare the raw password with the stored bcrypt hash
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(p))
	return err == nil
}

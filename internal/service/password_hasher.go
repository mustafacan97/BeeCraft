package services

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// PasswordHasher handles password hashing and verification.
type PasswordHasher struct{}

// HashWithBcrypt hashes the password using bcrypt (recommended for security).
// Built-in salting mechanism
func (ph *PasswordHasher) HashWithBcrypt(password string) (string, error) {
	// Generate bcrypt hash of the password with a default cost factor (10)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password with bcrypt: %v", err)
	}
	return string(hashedPassword), nil
}

// VerifyBcryptPassword verifies if the given raw password matches the bcrypt hash.
func (ph *PasswordHasher) VerifyBcryptPassword(rawPassword, storedHash string) (bool, error) {
	// Compare the raw password with the stored bcrypt hash
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(rawPassword))
	if err != nil {
		return false, fmt.Errorf("error comparing bcrypt hash: %v", err)
	}
	return true, nil
}

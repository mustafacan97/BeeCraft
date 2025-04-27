package valueobject

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"platform/pkg/domain"
)

type TraditionalCredential struct {
	domain.BaseValueObject
	username string
	password string
}

func NewTraditionalCredentials(username, password string) *TraditionalCredential {
	return &TraditionalCredential{
		username: username,
		password: password,
	}
}

func EncryptPassword(rawPassword string, secretKey []byte) (*string, error) {
	encrypted, err := encryptAES(rawPassword, secretKey)
	if err != nil {
		return nil, err
	}
	return &encrypted, nil
}

func DecryptPassword(password string, secretKey []byte) (*string, error) {
	decrypted, err := decryptAES(password, secretKey)
	if err != nil {
		return nil, err
	}
	return &decrypted, nil
}

func (e *TraditionalCredential) SetTraditionalCredentials(username, password string) {
	e.username = username
	e.password = password
}

func (e *TraditionalCredential) GetCredentials() (string, string) {
	return e.username, e.password
}

func (e *TraditionalCredential) GetAtomicValues() []interface{} {
	return []any{e.username, e.password}
}

func encryptAES(plainText string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func decryptAES(encryptedText string, key []byte) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("cipher text too short")
	}

	nonce, cipherText := data[:nonceSize], data[nonceSize:]
	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}
	return string(plainText), nil
}

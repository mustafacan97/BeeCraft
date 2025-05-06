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

var aesKey = []byte("1234567890123456")

type TraditionalCredential struct {
	domain.BaseValueObject
	username string
	password string
}

func NewTraditionalCredentials(username, rawPassword string) (*TraditionalCredential, error) {
	if username == "" || rawPassword == "" {
		return nil, errors.New("username and password are request for traditional authentication")
	}
	return &TraditionalCredential{
		username: username,
		password: rawPassword,
	}, nil
}

func (e *TraditionalCredential) GetCredentials() (string, string) {
	return e.username, e.password
}

func (e *TraditionalCredential) SetCredentials(username, password string) {
	e.username = username
	e.password = password
}

func (e *TraditionalCredential) GetAtomicValues() []interface{} {
	return []any{e.username, e.password}
}

func EncryptAES(plainText string) (string, error) {
	block, err := aes.NewCipher(aesKey)
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

func DecryptAES(encryptedText string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(aesKey)
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

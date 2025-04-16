package valueobject

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"platform/pkg/domain"

	"go.uber.org/zap"
)

type TraditionalCredential struct {
	domain.BaseValueObject
	username string
	password string
}

func NewTraditionalCredential(username, rawPassword string, secretKey []byte) *TraditionalCredential {
	encrypted, err := encryptAES(rawPassword, secretKey)
	if err != nil {
		zap.L().Error("failed to encrypt password",
			zap.String("username", username),
			zap.Error(err),
		)
		return nil
	}
	return &TraditionalCredential{
		username: username,
		password: encrypted,
	}
}

func (e *TraditionalCredential) GetCredentials(secretKey []byte) (string, string) {
	decrypted, err := decryptAES(e.password, secretKey)
	if err != nil {
		zap.L().Error("failed to decrypt password",
			zap.String("username", e.username),
			zap.Error(err),
		)
		return "", ""
	}
	return e.username, decrypted
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

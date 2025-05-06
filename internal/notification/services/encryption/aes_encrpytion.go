package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

type AESEncryptionService struct {
	key []byte
}

func NewAESEncryptionService(key []byte) (*AESEncryptionService, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, errors.New("invalid AES key length: must be 16, 24, or 32 bytes")
	}
	return &AESEncryptionService{key: key}, nil
}

func (s *AESEncryptionService) Encrypt(plainText string) (string, error) {
	if plainText == "" {
		return "", nil
	}

	block, err := aes.NewCipher(s.key)
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

func (s *AESEncryptionService) Decrypt(cipherText string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(s.key)
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
	nonce, cipherData := data[:nonceSize], data[nonceSize:]
	plainText, err := aesGCM.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return "", err
	}
	return string(plainText), nil
}

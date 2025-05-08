package encryption

type EncryptionService interface {
	Encrypt(plainText string) (string, error)
	Decrypt(cipherText string) (string, error)
}

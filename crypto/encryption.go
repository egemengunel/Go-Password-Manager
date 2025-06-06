package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	keySize   = 32 // AES-256
	nonceSize = 12 // GCM standard nonce size
	saltSize  = 32 // Salt size for key derivation
)

// EncryptionKey represents a derived encryption key
type EncryptionKey struct {
	Key  []byte
	Salt []byte
}

// DeriveKey derives an encryption key from a master password using PBKDF2
func DeriveKey(masterPassword string, salt []byte) *EncryptionKey {
	if salt == nil {
		salt = make([]byte, saltSize)
		if _, err := io.ReadFull(rand.Reader, salt); err != nil {
			panic(err) // This should be handled properly in production
		}
	}

	key := pbkdf2.Key([]byte(masterPassword), salt, 10000, keySize, sha256.New)

	return &EncryptionKey{
		Key:  key,
		Salt: salt,
	}
}

// Encrypt encrypts plaintext using AES-GCM
func Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	if len(key) != keySize {
		return nil, errors.New("invalid key size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt decrypts ciphertext using AES-GCM
func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	if len(key) != keySize {
		return nil, errors.New("invalid key size")
	}

	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := ciphertext[:nonceSize]
	ciphertext = ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// SecureZero clears sensitive data from memory
func SecureZero(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

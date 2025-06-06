package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/egemengunel/Go-Password-Manager/crypto"
	"github.com/egemengunel/Go-Password-Manager/models"
)

const VaultVersion = "1.0.0"

// VaultFile represents the encrypted vault file structure
type VaultFile struct {
	Version       string `json:"version"`
	PasswordHash  string `json:"password_hash"`
	Salt          []byte `json:"salt"`
	EncryptedData []byte `json:"encrypted_data"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// CreateVault creates a new encrypted vault file
func CreateVault(masterPassword, path string) error {
	// Check if vault already exists
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("vault already exists at %s", path)
	}

	// Hash the master password
	passwordHash, err := HashMasterPassword(masterPassword)
	if err != nil {
		return fmt.Errorf("failed to hash master password: %w", err)
	}

	// Derive encryption key
	encKey := crypto.DeriveKey(masterPassword, nil)

	// Create new vault structure
	vault := &models.Vault{
		Version:   VaultVersion,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Salt:      encKey.Salt,
		Entries:   make(map[string]*models.Entry),
		Metadata:  make(map[string]string),
	}

	// Save the vault
	return SaveVault(vault, path, encKey.Key, passwordHash)
}

// OpenVault opens and decrypts an existing vault
func OpenVault(masterPassword, path string) (*models.Vault, error) {
	// Check if vault exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("vault not found at %s", path)
	}

	// Read vault file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read vault file: %w", err)
	}

	// Parse vault file
	var vaultFile VaultFile
	if err := json.Unmarshal(data, &vaultFile); err != nil {
		return nil, fmt.Errorf("failed to parse vault file: %w", err)
	}

	// Verify master password
	match, err := VerifyMasterPassword(masterPassword, vaultFile.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("failed to verify password: %w", err)
	}
	if !match {
		return nil, fmt.Errorf("invalid master password")
	}

	// Derive decryption key
	encKey := crypto.DeriveKey(masterPassword, vaultFile.Salt)

	// Decrypt vault data
	decryptedData, err := crypto.Decrypt(vaultFile.EncryptedData, encKey.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt vault: %w", err)
	}

	// Parse decrypted vault
	var vault models.Vault
	if err := json.Unmarshal(decryptedData, &vault); err != nil {
		return nil, fmt.Errorf("failed to parse decrypted vault: %w", err)
	}

	return &vault, nil
}

// SaveVault encrypts and saves a vault to disk
func SaveVault(vault *models.Vault, path string, key []byte, passwordHash string) error {
	// Update timestamp
	vault.UpdatedAt = time.Now()

	// Marshal vault to JSON
	vaultData, err := json.Marshal(vault)
	if err != nil {
		return fmt.Errorf("failed to marshal vault: %w", err)
	}

	// Encrypt vault data
	encryptedData, err := crypto.Encrypt(vaultData, key)
	if err != nil {
		return fmt.Errorf("failed to encrypt vault: %w", err)
	}

	// Create vault file structure
	vaultFile := VaultFile{
		Version:       VaultVersion,
		PasswordHash:  passwordHash,
		Salt:          vault.Salt,
		EncryptedData: encryptedData,
		CreatedAt:     vault.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     vault.UpdatedAt.Format(time.RFC3339),
	}

	// Marshal vault file
	fileData, err := json.MarshalIndent(vaultFile, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal vault file: %w", err)
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("failed to create vault directory: %w", err)
	}

	// Write to file with secure permissions
	if err := os.WriteFile(path, fileData, 0600); err != nil {
		return fmt.Errorf("failed to write vault file: %w", err)
	}

	// Clear sensitive data from memory
	crypto.SecureZero(vaultData)
	crypto.SecureZero(key)

	return nil
}

// VaultExists checks if a vault file exists at the given path
func VaultExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

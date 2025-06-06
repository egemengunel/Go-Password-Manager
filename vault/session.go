package vault

import (
	"fmt"
	"sync"
	"time"

	"github.com/egemengunel/Go-Password-Manager/crypto"
	"github.com/egemengunel/Go-Password-Manager/models"
)

// Session represents an active vault session
type Session struct {
	Vault          *models.Vault
	VaultPath      string
	EncryptionKey  []byte
	PasswordHash   string
	LastAccessed   time.Time
	SessionTimeout time.Duration
	mutex          sync.RWMutex
}

var (
	currentSession *Session
	sessionMutex   sync.RWMutex
)

// StartSession creates a new vault session
func StartSession(vault *models.Vault, vaultPath string, masterPassword string, passwordHash string) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	// Derive encryption key
	encKey := crypto.DeriveKey(masterPassword, vault.Salt)

	currentSession = &Session{
		Vault:          vault,
		VaultPath:      vaultPath,
		EncryptionKey:  encKey.Key,
		PasswordHash:   passwordHash,
		LastAccessed:   time.Now(),
		SessionTimeout: 15 * time.Minute, // Default 15 minute timeout
	}
}

// GetSession returns the current active session
func GetSession() *Session {
	sessionMutex.RLock()
	defer sessionMutex.RUnlock()

	if currentSession == nil {
		return nil
	}

	// Check if session has expired
	if time.Since(currentSession.LastAccessed) > currentSession.SessionTimeout {
		// Session expired, clear it
		ClearSession()
		return nil
	}

	// Update last accessed time
	currentSession.LastAccessed = time.Now()
	return currentSession
}

// ClearSession clears the current session and sensitive data
func ClearSession() {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	if currentSession != nil {
		// Clear sensitive data from memory
		crypto.SecureZero(currentSession.EncryptionKey)
		currentSession = nil
	}
}

// IsSessionActive checks if there's an active session
func IsSessionActive() bool {
	return GetSession() != nil
}

// SaveCurrentSession saves the current session's vault to disk
func SaveCurrentSession() error {
	session := GetSession()
	if session == nil {
		return fmt.Errorf("no active session")
	}

	session.mutex.Lock()
	defer session.mutex.Unlock()

	return SaveVault(session.Vault, session.VaultPath, session.EncryptionKey, session.PasswordHash)
}

// AddEntryToSession adds an entry to the current session
func (s *Session) AddEntry(entry *models.Entry) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.Vault.Entries == nil {
		s.Vault.Entries = make(map[string]*models.Entry)
	}

	s.Vault.Entries[entry.ID] = entry
	s.LastAccessed = time.Now()
	return nil
}

// GetEntryFromSession retrieves an entry from the current session
func (s *Session) GetEntry(id string) (*models.Entry, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	entry, exists := s.Vault.Entries[id]
	if !exists {
		return nil, fmt.Errorf("entry not found")
	}

	s.LastAccessed = time.Now()
	return entry, nil
}

// ListEntriesFromSession returns all entries from the current session
func (s *Session) ListEntries() []*models.Entry {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	entries := make([]*models.Entry, 0, len(s.Vault.Entries))
	for _, entry := range s.Vault.Entries {
		entries = append(entries, entry)
	}

	s.LastAccessed = time.Now()
	return entries
}

// UpdateEntryInSession updates an entry in the current session
func (s *Session) UpdateEntry(entry *models.Entry) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.Vault.Entries[entry.ID]; !exists {
		return fmt.Errorf("entry not found")
	}

	entry.UpdatedAt = time.Now()
	s.Vault.Entries[entry.ID] = entry
	s.LastAccessed = time.Now()
	return nil
}

// DeleteEntryFromSession removes an entry from the current session
func (s *Session) DeleteEntry(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.Vault.Entries[id]; !exists {
		return fmt.Errorf("entry not found")
	}

	delete(s.Vault.Entries, id)
	s.LastAccessed = time.Now()
	return nil
}

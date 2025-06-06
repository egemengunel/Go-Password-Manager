package models

import (
	"time"
)

// Entry represents a password entry in the vault
type Entry struct {
	ID         string            `json:"id"`
	Title      string            `json:"title"`
	Username   string            `json:"username"`
	Password   string            `json:"password"`
	URL        string            `json:"url,omitempty"`
	Notes      string            `json:"notes,omitempty"`
	Tags       []string          `json:"tags,omitempty"`
	Custom     map[string]string `json:"custom,omitempty"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	AccessedAt time.Time         `json:"accessed_at"`
}

// Vault represents the structure of the password vault
type Vault struct {
	Version   string            `json:"version"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Salt      []byte            `json:"salt"`
	Entries   map[string]*Entry `json:"entries"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// NewEntry creates a new password entry with generated ID and timestamps
func NewEntry(title, username, password string) *Entry {
	now := time.Now()
	return &Entry{
		ID:         generateID(),
		Title:      title,
		Username:   username,
		Password:   password,
		CreatedAt:  now,
		UpdatedAt:  now,
		AccessedAt: now,
		Tags:       make([]string, 0),
		Custom:     make(map[string]string),
	}
}

// generateID creates a unique identifier for entries
func generateID() string {
	// This is a placeholder - should use crypto/rand for production
	return time.Now().Format("20060102150405") + "_" + "entry"
}

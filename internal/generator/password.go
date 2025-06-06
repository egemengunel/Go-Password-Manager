package generator

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const (
	lowercase = "abcdefghijklmnopqrstuvwxyz"
	uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers   = "0123456789"
	symbols   = "!@#$%^&*()_+-=[]{}|;:,.<>?"
)

// PasswordOptions defines options for password generation
type PasswordOptions struct {
	Length           int
	IncludeLower     bool
	IncludeUpper     bool
	IncludeNumbers   bool
	IncludeSymbols   bool
	ExcludeAmbiguous bool
}

// DefaultOptions returns sensible default options for password generation
func DefaultOptions() PasswordOptions {
	return PasswordOptions{
		Length:           16,
		IncludeLower:     true,
		IncludeUpper:     true,
		IncludeNumbers:   true,
		IncludeSymbols:   true,
		ExcludeAmbiguous: true,
	}
}

// GeneratePassword creates a secure random password based on the given options
func GeneratePassword(opts PasswordOptions) (string, error) {
	if opts.Length <= 0 {
		opts.Length = 16
	}

	var charset string

	if opts.IncludeLower {
		charset += lowercase
	}
	if opts.IncludeUpper {
		charset += uppercase
	}
	if opts.IncludeNumbers {
		charset += numbers
	}
	if opts.IncludeSymbols {
		charset += symbols
	}

	if opts.ExcludeAmbiguous {
		// Remove ambiguous characters
		ambiguous := "0O1lI|"
		for _, char := range ambiguous {
			charset = strings.ReplaceAll(charset, string(char), "")
		}
	}

	if len(charset) == 0 {
		charset = lowercase + uppercase + numbers // fallback
	}

	password := make([]byte, opts.Length)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := range password {
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		password[i] = charset[randomIndex.Int64()]
	}

	return string(password), nil
}

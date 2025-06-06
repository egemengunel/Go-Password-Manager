package display

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"

	"github.com/egemengunel/Go-Password-Manager/models"
)

// Colors for different elements
var (
	successColor = color.New(color.FgGreen, color.Bold)
	errorColor   = color.New(color.FgRed, color.Bold)
	infoColor    = color.New(color.FgCyan)
	warningColor = color.New(color.FgYellow)
	titleColor   = color.New(color.FgMagenta, color.Bold)
)

// Success prints a success message
func Success(message string) {
	successColor.Printf("✓ %s\n", message)
}

// Error prints an error message
func Error(message string) {
	errorColor.Printf("✗ %s\n", message)
}

// Info prints an info message
func Info(message string) {
	infoColor.Printf("ℹ %s\n", message)
}

// Warning prints a warning message
func Warning(message string) {
	warningColor.Printf("⚠ %s\n", message)
}

// Title prints a section title
func Title(message string) {
	titleColor.Printf("\n%s\n", message)
	titleColor.Printf("%s\n\n", strings.Repeat("=", len(message)))
}

// MaskPassword masks a password for display
func MaskPassword(password string) string {
	if len(password) == 0 {
		return ""
	}
	return strings.Repeat("*", min(len(password), 8))
}

// FormatTime formats a time for display
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// FormatTimeAgo formats time as "X ago"
func FormatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	if duration < time.Minute {
		return "just now"
	} else if duration < time.Hour {
		minutes := int(duration.Minutes())
		return fmt.Sprintf("%d min ago", minutes)
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		return fmt.Sprintf("%d hr ago", hours)
	} else {
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%d days ago", days)
	}
}

// ListEntries displays entries in a table format
func ListEntries(entries []*models.Entry, showPasswords bool) {
	if len(entries) == 0 {
		Info("No entries found")
		return
	}

	// Sort entries by title
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Title < entries[j].Title
	})

	// Simple table formatting without complex tablewriter features
	fmt.Printf("%-3s %-20s %-20s %-15s %-30s %-15s\n",
		"#", "Title", "Username", "Password", "URL", "Updated")
	fmt.Printf("%s\n", strings.Repeat("-", 100))

	for i, entry := range entries {
		password := MaskPassword(entry.Password)
		if showPasswords {
			password = entry.Password
		}

		url := entry.URL
		if len(url) > 28 {
			url = url[:25] + "..."
		}

		title := entry.Title
		if len(title) > 18 {
			title = title[:15] + "..."
		}

		username := entry.Username
		if len(username) > 18 {
			username = username[:15] + "..."
		}

		fmt.Printf("%-3d %-20s %-20s %-15s %-30s %-15s\n",
			i+1, title, username, password, url, FormatTimeAgo(entry.UpdatedAt))
	}

	fmt.Printf("\nTotal: %d entries\n", len(entries))
}

// ShowEntryDetails displays detailed information about a single entry
func ShowEntryDetails(entry *models.Entry, showPassword bool) {
	Title(fmt.Sprintf("Entry: %s", entry.Title))

	fmt.Printf("ID:         %s\n", entry.ID)
	fmt.Printf("Title:      %s\n", entry.Title)
	fmt.Printf("Username:   %s\n", entry.Username)

	if showPassword {
		fmt.Printf("Password:   %s\n", entry.Password)
	} else {
		fmt.Printf("Password:   %s\n", MaskPassword(entry.Password))
	}

	if entry.URL != "" {
		fmt.Printf("URL:        %s\n", entry.URL)
	}

	if entry.Notes != "" {
		fmt.Printf("Notes:      %s\n", entry.Notes)
	}

	if len(entry.Tags) > 0 {
		fmt.Printf("Tags:       %s\n", strings.Join(entry.Tags, ", "))
	}

	if len(entry.Custom) > 0 {
		fmt.Printf("Custom:     ")
		for key, value := range entry.Custom {
			fmt.Printf("%s=%s ", key, value)
		}
		fmt.Println()
	}

	fmt.Printf("Created:    %s\n", FormatTime(entry.CreatedAt))
	fmt.Printf("Updated:    %s\n", FormatTime(entry.UpdatedAt))
	fmt.Printf("Accessed:   %s\n", FormatTime(entry.AccessedAt))
}

// ConfirmAction prompts for confirmation before dangerous actions
func ConfirmAction(action, target string) bool {
	warningColor.Printf("⚠ Are you sure you want to %s '%s'? This action cannot be undone.\n", action, target)
	fmt.Print("Type 'yes' to confirm: ")

	var response string
	fmt.Scanln(&response)

	return strings.ToLower(response) == "yes"
}

// ShowVaultInfo displays information about the vault
func ShowVaultInfo(vault *models.Vault) {
	Title("Vault Information")

	fmt.Printf("Version:    %s\n", vault.Version)
	fmt.Printf("Created:    %s\n", FormatTime(vault.CreatedAt))
	fmt.Printf("Updated:    %s\n", FormatTime(vault.UpdatedAt))
	fmt.Printf("Entries:    %d\n", len(vault.Entries))

	if len(vault.Metadata) > 0 {
		fmt.Printf("Metadata:   ")
		for key, value := range vault.Metadata {
			fmt.Printf("%s=%s ", key, value)
		}
		fmt.Println()
	}
}

// ShowPasswordStrength displays password strength information
func ShowPasswordStrength(password string) {
	length := len(password)
	hasUpper := strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	hasLower := strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz")
	hasDigit := strings.ContainsAny(password, "0123456789")
	hasSymbol := strings.ContainsAny(password, "!@#$%^&*()_+-=[]{}|;:,.<>?")

	score := 0
	if length >= 8 {
		score++
	}
	if length >= 12 {
		score++
	}
	if hasUpper {
		score++
	}
	if hasLower {
		score++
	}
	if hasDigit {
		score++
	}
	if hasSymbol {
		score++
	}

	var strength string
	var strengthColor *color.Color

	switch {
	case score >= 5:
		strength = "Strong"
		strengthColor = color.New(color.FgGreen)
	case score >= 3:
		strength = "Medium"
		strengthColor = color.New(color.FgYellow)
	default:
		strength = "Weak"
		strengthColor = color.New(color.FgRed)
	}

	fmt.Printf("Password strength: ")
	strengthColor.Printf("%s", strength)
	fmt.Printf(" (Length: %d, Score: %d/6)\n", length, score)
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

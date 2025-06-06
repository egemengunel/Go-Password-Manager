package cmd

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/egemengunel/Go-Password-Manager/config"
	"github.com/egemengunel/Go-Password-Manager/internal/display"
	"github.com/egemengunel/Go-Password-Manager/internal/generator"
	"github.com/egemengunel/Go-Password-Manager/internal/input"
	"github.com/egemengunel/Go-Password-Manager/models"
	"github.com/egemengunel/Go-Password-Manager/vault"
)

var editCmd = &cobra.Command{
	Use:   "edit <entry-id-or-number>",
	Short: "Edit an existing password entry",
	Long: `Edit an existing password entry in your vault.
You can specify the entry by ID or by its number from the list command.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runEdit(cmd, args)
	},
}

var (
	editTitle    string
	editUsername string
	editPassword string
	editURL      string
	editNotes    string
	editGenerate bool
	editLength   int
)

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.Flags().StringVarP(&editTitle, "title", "t", "", "New title for the entry")
	editCmd.Flags().StringVarP(&editUsername, "username", "u", "", "New username for the entry")
	editCmd.Flags().StringVarP(&editPassword, "password", "p", "", "New password for the entry")
	editCmd.Flags().StringVar(&editURL, "url", "", "New URL for the entry")
	editCmd.Flags().StringVar(&editNotes, "notes", "", "New notes for the entry")
	editCmd.Flags().BoolVarP(&editGenerate, "generate", "g", false, "Generate a new random password")
	editCmd.Flags().IntVarP(&editLength, "length", "l", 16, "Length of generated password")
}

func runEdit(cmd *cobra.Command, args []string) {
	// Get configuration
	cfg := config.DefaultConfig()

	// Check if vault exists
	if !vault.VaultExists(cfg.VaultPath) {
		display.Error("No vault found. Please run 'gopassman init' first")
		os.Exit(1)
	}

	// Try to use existing session first
	session := vault.GetSession()
	if session == nil {
		// No active session, need to open vault
		if !input.CheckTTY() {
			display.Error("No active session. Please run with a valid session or in interactive mode")
			os.Exit(1)
		}

		// Prompt for master password
		masterPassword, err := input.PromptMasterPassword("Enter master password: ")
		if err != nil {
			display.Error(fmt.Sprintf("Failed to read password: %v", err))
			os.Exit(1)
		}

		// Open vault
		vaultData, err := vault.OpenVault(masterPassword, cfg.VaultPath)
		if err != nil {
			display.Error(fmt.Sprintf("Failed to open vault: %v", err))
			os.Exit(1)
		}

		// Get password hash for session
		passwordHash, err := vault.HashMasterPassword(masterPassword)
		if err != nil {
			display.Error(fmt.Sprintf("Failed to hash password: %v", err))
			os.Exit(1)
		}

		// Start session
		vault.StartSession(vaultData, cfg.VaultPath, masterPassword, passwordHash)
		session = vault.GetSession()

		// Clear master password from memory
		for i := range masterPassword {
			masterPassword = masterPassword[:i] + "x" + masterPassword[i+1:]
		}
	}

	entryIdentifier := args[0]
	var entry *models.Entry
	var err error

	// Try to find entry by ID first
	entry, err = session.GetEntry(entryIdentifier)
	if err != nil {
		// If not found by ID, try by number
		if num, parseErr := strconv.Atoi(entryIdentifier); parseErr == nil {
			entries := session.ListEntries()
			if num > 0 && num <= len(entries) {
				// Sort entries by title (same as list command)
				sortedEntries := make([]*models.Entry, len(entries))
				copy(sortedEntries, entries)
				sort.Slice(sortedEntries, func(i, j int) bool {
					return sortedEntries[i].Title < sortedEntries[j].Title
				})

				entry = sortedEntries[num-1]
			} else {
				display.Error(fmt.Sprintf("Entry number %d not found. Use 'gopassman list' to see available entries", num))
				os.Exit(1)
			}
		} else {
			display.Error(fmt.Sprintf("Entry '%s' not found", entryIdentifier))
			os.Exit(1)
		}
	}

	// Show current entry details
	fmt.Printf("Editing entry: %s\n", entry.Title)
	display.ShowEntryDetails(entry, false)

	// Check if any flags were provided
	hasFlags := editTitle != "" || editUsername != "" || editPassword != "" ||
		editURL != "" || editNotes != "" || editGenerate

	if hasFlags {
		// Use flag values
		if editTitle != "" {
			entry.Title = editTitle
		}
		if editUsername != "" {
			entry.Username = editUsername
		}
		if editPassword != "" {
			entry.Password = editPassword
		}
		if editURL != "" {
			entry.URL = editURL
		}
		if editNotes != "" {
			entry.Notes = editNotes
		}

		// Generate password if requested
		if editGenerate {
			opts := generator.DefaultOptions()
			opts.Length = editLength
			if editLength < 8 {
				opts.Length = 16
			}

			generatedPassword, err := generator.GeneratePassword(opts)
			if err != nil {
				display.Error(fmt.Sprintf("Failed to generate password: %v", err))
				os.Exit(1)
			}
			entry.Password = generatedPassword

			display.Info(fmt.Sprintf("Generated new password: %s", generatedPassword))
			display.ShowPasswordStrength(generatedPassword)
		}
	} else {
		// Interactive mode
		if !input.CheckTTY() {
			display.Error("Interactive mode requires a terminal. Use flags instead")
			os.Exit(1)
		}

		display.Title("Edit Password Entry")
		display.Info("Leave blank to keep current value")

		// Title
		if newTitle, err := input.PromptString(fmt.Sprintf("Title [%s]:", entry.Title), false); err == nil && newTitle != "" {
			entry.Title = newTitle
		}

		// Username
		if newUsername, err := input.PromptString(fmt.Sprintf("Username [%s]:", entry.Username), false); err == nil && newUsername != "" {
			entry.Username = newUsername
		}

		// Password
		generatePassword, err := input.PromptConfirm("Generate new password?", false)
		if err != nil {
			display.Error(fmt.Sprintf("Failed to get password choice: %v", err))
			os.Exit(1)
		}

		if generatePassword {
			opts := generator.DefaultOptions()
			generatedPassword, err := generator.GeneratePassword(opts)
			if err != nil {
				display.Error(fmt.Sprintf("Failed to generate password: %v", err))
				os.Exit(1)
			}
			entry.Password = generatedPassword

			display.Info(fmt.Sprintf("Generated new password: %s", generatedPassword))
			display.ShowPasswordStrength(generatedPassword)
		} else {
			if newPassword, err := input.PromptPassword(fmt.Sprintf("Password [%s]:", display.MaskPassword(entry.Password)), false); err == nil && newPassword != "" {
				entry.Password = newPassword
			}
		}

		// URL
		if newURL, err := input.PromptString(fmt.Sprintf("URL [%s]:", entry.URL), false); err == nil && newURL != "" {
			entry.URL = newURL
		}

		// Notes
		if newNotes, err := input.PromptMultiline(fmt.Sprintf("Notes [%s]:", entry.Notes)); err == nil && newNotes != "" {
			entry.Notes = newNotes
		}
	}

	// Update timestamps
	entry.UpdatedAt = time.Now()

	// Update entry in session
	if err := session.UpdateEntry(entry); err != nil {
		display.Error(fmt.Sprintf("Failed to update entry: %v", err))
		os.Exit(1)
	}

	// Save vault
	if err := vault.SaveCurrentSession(); err != nil {
		display.Error(fmt.Sprintf("Failed to save vault: %v", err))
		os.Exit(1)
	}

	display.Success(fmt.Sprintf("Entry '%s' updated successfully", entry.Title))
}

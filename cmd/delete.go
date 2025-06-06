package cmd

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/egemengunel/Go-Password-Manager/config"
	"github.com/egemengunel/Go-Password-Manager/internal/display"
	"github.com/egemengunel/Go-Password-Manager/internal/input"
	"github.com/egemengunel/Go-Password-Manager/models"
	"github.com/egemengunel/Go-Password-Manager/vault"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <entry-id-or-number>",
	Short: "Delete a password entry",
	Long: `Delete a password entry from your vault.
You can specify the entry by ID or by its number from the list command.
This action cannot be undone.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runDelete(cmd, args)
	},
}

var deleteForce bool

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVarP(&deleteForce, "force", "f", false, "Force deletion without confirmation")
}

func runDelete(cmd *cobra.Command, args []string) {
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

	// Show entry details before deletion
	display.Title(fmt.Sprintf("Delete Entry: %s", entry.Title))
	display.ShowEntryDetails(entry, false)

	// Confirm deletion unless forced
	if !deleteForce {
		if !input.CheckTTY() {
			display.Error("Deletion requires confirmation. Use --force to bypass or run in interactive mode")
			os.Exit(1)
		}

		confirmed, err := input.PromptConfirm(fmt.Sprintf("Are you sure you want to delete '%s'?", entry.Title), false)
		if err != nil {
			display.Error(fmt.Sprintf("Failed to get confirmation: %v", err))
			os.Exit(1)
		}

		if !confirmed {
			display.Info("Deletion cancelled")
			return
		}

		// Double confirmation for safety
		display.Warning("This action cannot be undone!")
		doubleConfirmed, err := input.PromptConfirm("Type 'yes' to confirm deletion", false)
		if err != nil {
			display.Error(fmt.Sprintf("Failed to get confirmation: %v", err))
			os.Exit(1)
		}

		if !doubleConfirmed {
			display.Info("Deletion cancelled")
			return
		}
	}

	// Delete entry from session
	if err := session.DeleteEntry(entry.ID); err != nil {
		display.Error(fmt.Sprintf("Failed to delete entry: %v", err))
		os.Exit(1)
	}

	// Save vault
	if err := vault.SaveCurrentSession(); err != nil {
		display.Error(fmt.Sprintf("Failed to save vault: %v", err))
		os.Exit(1)
	}

	display.Success(fmt.Sprintf("Entry '%s' deleted successfully", entry.Title))
	display.Info("Entry has been permanently removed from your vault")
}

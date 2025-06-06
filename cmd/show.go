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
	"github.com/egemengunel/Go-Password-Manager/internal/input"
	"github.com/egemengunel/Go-Password-Manager/models"
	"github.com/egemengunel/Go-Password-Manager/vault"
)

var showCmd = &cobra.Command{
	Use:   "show <entry-id-or-number>",
	Short: "Show detailed information about an entry",
	Long: `Show detailed information about a specific password entry.
You can specify the entry by ID or by its number from the list command.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runShow(cmd, args)
	},
}

var showPassword bool

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.Flags().BoolVarP(&showPassword, "password", "p", false, "Show the password in plain text")
}

func runShow(cmd *cobra.Command, args []string) {
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

	// Update access time
	entry.AccessedAt = time.Now()
	if err := session.UpdateEntry(entry); err != nil {
		display.Warning("Failed to update access time")
	}

	// Save the updated access time
	if err := vault.SaveCurrentSession(); err != nil {
		display.Warning("Failed to save access time update")
	}

	// Display entry details
	display.ShowEntryDetails(entry, showPassword)

	if !showPassword {
		fmt.Println()
		display.Info("Use --password to show the password in plain text")
	}
}

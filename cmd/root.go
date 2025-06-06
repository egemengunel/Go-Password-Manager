package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gopassman",
	Short: "A secure CLI password manager",
	Long: `Go Password Manager is a secure command-line interface for managing your passwords.
It uses strong encryption (AES-GCM) and secure key derivation (Argon2id) to protect your data.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Global flags can be added here
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gopassman.yaml)")
}

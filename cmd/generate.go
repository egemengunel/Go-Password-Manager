package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/egemengunel/Go-Password-Manager/internal/display"
	"github.com/egemengunel/Go-Password-Manager/internal/generator"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a secure password",
	Long: `Generate a secure random password with customizable options.
This command does not store the password - it only generates and displays it.`,
	Run: func(cmd *cobra.Command, args []string) {
		runGenerate(cmd, args)
	},
}

var (
	genLength    int
	genNoUpper   bool
	genNoLower   bool
	genNoNumbers bool
	genNoSymbols bool
	genAmbiguous bool
	genCount     int
)

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().IntVarP(&genLength, "length", "l", 16, "Length of the password")
	generateCmd.Flags().BoolVar(&genNoUpper, "no-upper", false, "Exclude uppercase letters")
	generateCmd.Flags().BoolVar(&genNoLower, "no-lower", false, "Exclude lowercase letters")
	generateCmd.Flags().BoolVar(&genNoNumbers, "no-numbers", false, "Exclude numbers")
	generateCmd.Flags().BoolVar(&genNoSymbols, "no-symbols", false, "Exclude symbols")
	generateCmd.Flags().BoolVar(&genAmbiguous, "ambiguous", false, "Include ambiguous characters (0, O, 1, l, I, |)")
	generateCmd.Flags().IntVarP(&genCount, "count", "c", 1, "Number of passwords to generate")
}

func runGenerate(cmd *cobra.Command, args []string) {
	// Validate parameters
	if genLength < 4 {
		display.Error("Password length must be at least 4 characters")
		return
	}

	if genCount < 1 || genCount > 50 {
		display.Error("Count must be between 1 and 50")
		return
	}

	// Build password options
	opts := generator.PasswordOptions{
		Length:           genLength,
		IncludeLower:     !genNoLower,
		IncludeUpper:     !genNoUpper,
		IncludeNumbers:   !genNoNumbers,
		IncludeSymbols:   !genNoSymbols,
		ExcludeAmbiguous: !genAmbiguous,
	}

	// Validate that at least one character type is included
	if !opts.IncludeLower && !opts.IncludeUpper && !opts.IncludeNumbers && !opts.IncludeSymbols {
		display.Error("At least one character type must be included")
		return
	}

	display.Title("Generated Passwords")

	// Generate and display passwords
	for i := 0; i < genCount; i++ {
		password, err := generator.GeneratePassword(opts)
		if err != nil {
			display.Error(fmt.Sprintf("Failed to generate password: %v", err))
			return
		}

		if genCount == 1 {
			fmt.Printf("Password: %s\n", password)
			display.ShowPasswordStrength(password)
		} else {
			fmt.Printf("%2d: %s\n", i+1, password)
		}
	}

	// Show generation settings
	fmt.Println()
	display.Info(fmt.Sprintf("Settings: Length=%d", genLength))

	var included []string
	if opts.IncludeLower {
		included = append(included, "lowercase")
	}
	if opts.IncludeUpper {
		included = append(included, "uppercase")
	}
	if opts.IncludeNumbers {
		included = append(included, "numbers")
	}
	if opts.IncludeSymbols {
		included = append(included, "symbols")
	}

	if len(included) > 0 {
		display.Info(fmt.Sprintf("Includes: %s", fmt.Sprintf("%v", included)))
	}

	if opts.ExcludeAmbiguous {
		display.Info("Ambiguous characters excluded")
	}
}

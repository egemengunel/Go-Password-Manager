package input

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/AlecAivazis/survey/v2"
	"golang.org/x/term"
)

// PromptMasterPassword securely prompts for master password
func PromptMasterPassword(message string) (string, error) {
	if message == "" {
		message = "Enter master password: "
	}

	fmt.Print(message)
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // Add newline after password input

	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}

	password := string(passwordBytes)
	if len(password) == 0 {
		return "", fmt.Errorf("password cannot be empty")
	}

	return password, nil
}

// PromptConfirmPassword prompts for password confirmation
func PromptConfirmPassword(original string) error {
	confirmation, err := PromptMasterPassword("Confirm master password: ")
	if err != nil {
		return err
	}

	if original != confirmation {
		return fmt.Errorf("passwords do not match")
	}

	return nil
}

// PromptString prompts for a string input
func PromptString(message string, required bool) (string, error) {
	prompt := &survey.Input{
		Message: message,
	}

	var result string
	err := survey.AskOne(prompt, &result)
	if err != nil {
		return "", err
	}

	result = strings.TrimSpace(result)
	if required && result == "" {
		return "", fmt.Errorf("input is required")
	}

	return result, nil
}

// PromptPassword prompts for a regular password (not master password)
func PromptPassword(message string, required bool) (string, error) {
	prompt := &survey.Password{
		Message: message,
	}

	var result string
	err := survey.AskOne(prompt, &result)
	if err != nil {
		return "", err
	}

	if required && result == "" {
		return "", fmt.Errorf("password is required")
	}

	return result, nil
}

// PromptConfirm prompts for yes/no confirmation
func PromptConfirm(message string, defaultValue bool) (bool, error) {
	prompt := &survey.Confirm{
		Message: message,
		Default: defaultValue,
	}

	var result bool
	err := survey.AskOne(prompt, &result)
	return result, err
}

// PromptSelect prompts for selection from options
func PromptSelect(message string, options []string) (string, error) {
	prompt := &survey.Select{
		Message: message,
		Options: options,
	}

	var result string
	err := survey.AskOne(prompt, &result)
	return result, err
}

// PromptMultiline prompts for multiline input
func PromptMultiline(message string) (string, error) {
	prompt := &survey.Multiline{
		Message: message,
	}

	var result string
	err := survey.AskOne(prompt, &result)
	return result, err
}

// PromptEntryDetails prompts for complete entry details
func PromptEntryDetails() (title, username, password, url, notes string, err error) {
	// Title (required)
	title, err = PromptString("Title:", true)
	if err != nil {
		return
	}

	// Username (required)
	username, err = PromptString("Username:", true)
	if err != nil {
		return
	}

	// Password (optional - can be generated)
	generatePassword, err := PromptConfirm("Generate password?", true)
	if err != nil {
		return
	}

	if generatePassword {
		password = "[GENERATED]" // Placeholder - will be replaced with generated password
	} else {
		password, err = PromptPassword("Password:", true)
		if err != nil {
			return
		}
	}

	// URL (optional)
	url, err = PromptString("URL (optional):", false)
	if err != nil {
		return
	}

	// Notes (optional)
	notes, err = PromptMultiline("Notes (optional):")
	if err != nil {
		return
	}

	return
}

// CheckTTY checks if we're running in an interactive terminal
func CheckTTY() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

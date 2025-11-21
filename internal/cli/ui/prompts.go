package ui

import (
	"bufio"
	"fmt"
	"strings"
)

// ConfirmationPrompt displays a yes/no confirmation prompt
// Returns true if user confirms (y/yes), false otherwise
// defaultNo: if true, defaults to "no" (safer option)
func ConfirmationPrompt(scanner *bufio.Scanner, message string, defaultNo bool) bool {
	defaultOption := "y/N"
	if !defaultNo {
		defaultOption = "Y/n"
	}

	fmt.Printf("%s (%s): ", message, defaultOption)

	if !scanner.Scan() {
		// If scan fails, return safe default
		return !defaultNo
	}

	input := strings.ToLower(strings.TrimSpace(scanner.Text()))

	// Handle empty input (use default)
	if input == "" {
		return !defaultNo
	}

	// Check for yes
	if input == "y" || input == "yes" {
		return true
	}

	// Check for no
	if input == "n" || input == "no" {
		return false
	}

	// Invalid input, ask again
	fmt.Println(Colorize("Invalid input. Please enter 'y' or 'n'.", ColorYellow))
	return ConfirmationPrompt(scanner, message, defaultNo)
}

// ConfirmDestructiveAction prompts for confirmation on destructive actions
// Always defaults to "no" for safety
func ConfirmDestructiveAction(scanner *bufio.Scanner, action string) bool {
	message := fmt.Sprintf("Are you sure you want to %s? This cannot be undone", action)
	return ConfirmationPrompt(scanner, message, true)
}

// ConfirmExpensivePurchase prompts for confirmation on expensive purchases
// Shows the cost and asks for confirmation
func ConfirmExpensivePurchase(scanner *bufio.Scanner, itemName string, cost int, currentCoins int) bool {
	fmt.Println()
	fmt.Printf("Purchase: %s\n", Colorize(itemName, Bold+ColorBrightCyan))
	fmt.Printf("Cost: %s coins\n", Colorize(fmt.Sprintf("%d", cost), ColorYellow))
	fmt.Printf("Your coins: %s\n", Colorize(fmt.Sprintf("%d", currentCoins), ColorGreen))
	fmt.Printf("After purchase: %s coins\n", Colorize(fmt.Sprintf("%d", currentCoins-cost), ColorGreen))
	fmt.Println()

	return ConfirmationPrompt(scanner, "Confirm purchase?", true)
}

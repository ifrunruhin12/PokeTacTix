// +build ignore

package main

import (
	"fmt"
	"os"
	"pokemon-cli/auth"
)

func main() {
	fmt.Println("=== Auth Package Verification ===\n")

	os.Setenv("JWT_SECRET", "test-secret-key-minimum-32-characters-required-for-security")

	fmt.Println("Test 1: Password Validation")
	service := auth.NewService()
	
	testPasswords := map[string]bool{
		"weak":                false, // Too short
		"NoSpecial123":        false, // No special char
		"nouppercas3!":        false, // No uppercase
		"NOLOWERCASE3!":       false, // No lowercase
		"NoNumbers!":          false, // No numbers
		"ValidPass123!":       true,  // Valid
		"Secure@Password99":   true,  // Valid
	}
	
	for password, shouldBeValid := range testPasswords {
		err := service.ValidatePassword(password)
		isValid := err == nil
		status := "✓"
		if isValid != shouldBeValid {
			status = "✗"
		}
		fmt.Printf("  %s Password '%s': valid=%v (expected=%v)\n", status, password, isValid, shouldBeValid)
	}

	fmt.Println("\nTest 2: Username Validation")
	testUsernames := map[string]bool{
		"ab":              false, // Too short
		"valid_user123":   true,  // Valid
		"user@name":       false, // Invalid chars
		"user name":       false, // Space
		"a":               false, // Too short
		"trainer_ash":     true,  // Valid
	}
	
	for username, shouldBeValid := range testUsernames {
		err := service.ValidateUsername(username)
		isValid := err == nil
		status := "✓"
		if isValid != shouldBeValid {
			status = "✗"
		}
		fmt.Printf("  %s Username '%s': valid=%v (expected=%v)\n", status, username, isValid, shouldBeValid)
	}

	fmt.Println("\nTest 3: Email Validation")
	testEmails := map[string]bool{
		"invalid":              false, // No @
		"test@example.com":     true,  // Valid
		"user@domain":          false, // No TLD
		"@example.com":         false, // No local part
		"ash.ketchum@poke.com": true,  // Valid
	}
	
	for email, shouldBeValid := range testEmails {
		err := service.ValidateEmail(email)
		isValid := err == nil
		status := "✓"
		if isValid != shouldBeValid {
			status = "✗"
		}
		fmt.Printf("  %s Email '%s': valid=%v (expected=%v)\n", status, email, isValid, shouldBeValid)
	}

	fmt.Println("\nTest 4: Password Hashing and Comparison")
	password := "SecurePass123!"
	hash, err := service.HashPassword(password)
	if err != nil {
		fmt.Printf("  ✗ Failed to hash password: %v\n", err)
	} else {
		fmt.Printf("  ✓ Password hashed successfully\n")

		err = service.ComparePassword(hash, password)
		if err == nil {
			fmt.Printf("  ✓ Correct password verified\n")
		} else {
			fmt.Printf("  ✗ Failed to verify correct password: %v\n", err)
		}

		err = service.ComparePassword(hash, "WrongPassword123!")
		if err != nil {
			fmt.Printf("  ✓ Wrong password correctly rejected\n")
		} else {
			fmt.Printf("  ✗ Wrong password was accepted (should fail)\n")
		}
	}

	fmt.Println("\nTest 5: JWT Token Generation and Validation")
	jwtService, err := auth.NewJWTService()
	if err != nil {
		fmt.Printf("  ✗ Failed to create JWT service: %v\n", err)
		return
	}
	
	token, err := jwtService.GenerateToken(123, "test_user")
	if err != nil {
		fmt.Printf("  ✗ Failed to generate token: %v\n", err)
	} else {
		fmt.Printf("  ✓ Token generated successfully\n")

		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			fmt.Printf("  ✗ Failed to validate token: %v\n", err)
		} else if claims.UserID == 123 && claims.Username == "test_user" {
			fmt.Printf("  ✓ Token validated successfully (UserID=%d, Username=%s)\n", claims.UserID, claims.Username)
		} else {
			fmt.Printf("  ✗ Token claims incorrect (UserID=%d, Username=%s)\n", claims.UserID, claims.Username)
		}

		_, err = jwtService.ValidateToken("invalid.token.here")
		if err != nil {
			fmt.Printf("  ✓ Invalid token correctly rejected\n")
		} else {
			fmt.Printf("  ✗ Invalid token was accepted (should fail)\n")
		}
	}
	
	fmt.Println("\n=== Verification Complete ===")
}

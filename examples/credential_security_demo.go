package main

import (
	"context"
	"fmt"

	"github.com/dshills/mcp-pr/internal/credentials"
	"github.com/dshills/mcp-pr/internal/logging"
)

// This example demonstrates the credential validation and masking features
func main() {
	// Initialize logging
	logging.Init("info")
	ctx := context.Background()

	// Create validator
	validator := credentials.NewValidator()

	// Example 1: Valid API keys
	fmt.Println("=== Example 1: Valid API Keys ===")
	anthropicKey := "sk-ant-api03-abcdefghijklmnopqrstuvwxyz1234567890"
	openaiKey := "sk-proj-abcdefghijklmnopqrstuvwxyz1234567890"

	if err := validator.ValidateAnthropicKey(anthropicKey); err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	} else {
		fmt.Println("✓ Anthropic key validated")
	}

	if err := validator.ValidateOpenAIKey(openaiKey); err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	} else {
		fmt.Println("✓ OpenAI key validated")
	}

	// Example 2: Demonstrate masking
	fmt.Println("\n=== Example 2: Key Masking ===")
	fmt.Printf("Original key: %s\n", anthropicKey)
	fmt.Printf("Masked key:   %s\n", credentials.MaskKey(anthropicKey))

	// Example 3: Logging with automatic masking
	fmt.Println("\n=== Example 3: Secure Logging ===")
	fmt.Println("Without masking (unsafe):")
	logging.Info(ctx, "Provider initialized",
		"provider", "anthropic",
		"api_key", anthropicKey,
	)

	fmt.Println("\nWith masking (safe):")
	logging.InfoWithMasking(ctx, "Provider initialized",
		"provider", "anthropic",
		"api_key", anthropicKey,
	)

	// Example 4: Invalid credentials
	fmt.Println("\n=== Example 4: Invalid Credentials ===")
	invalidKeys := []struct {
		name string
		key  string
	}{
		{"Empty key", ""},
		{"Too short", "sk-short"},
		{"Wrong prefix", "invalid-prefix-1234567890"},
		{"Placeholder", "your-api-key"},
	}

	for _, test := range invalidKeys {
		if err := validator.ValidateAnthropicKey(test.key); err != nil {
			fmt.Printf("✓ Correctly rejected '%s': %v\n", test.name, err)
		}
	}

	// Example 5: Multiple credential validation
	fmt.Println("\n=== Example 5: Batch Validation ===")
	err := validator.ValidateAll(
		"sk-ant-valid1234567890",
		"sk-proj-valid1234567890",
		"AIzaSyAbCdEfGhIjKlMnOpQrStUvWxYz1234567",
	)
	if err != nil {
		fmt.Printf("Validation errors: %v\n", err)
	} else {
		fmt.Println("✓ All credentials validated successfully")
	}
}

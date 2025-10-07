package credentials

import (
	"errors"
	"fmt"
	"strings"
)

// ValidationError represents a credential validation error
type ValidationError struct {
	Provider string
	Reason   string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid credential for %s: %s", e.Provider, e.Reason)
}

// Validator validates API credentials for different providers
type Validator struct {
	minLength int
	maxLength int
}

// NewValidator creates a new credential validator
func NewValidator() *Validator {
	return &Validator{
		minLength: 8,   // Minimum reasonable API key length
		maxLength: 512, // Maximum reasonable API key length
	}
}

// ValidateAnthropicKey validates an Anthropic API key
func (v *Validator) ValidateAnthropicKey(key string) error {
	if err := v.validateBasic("Anthropic", key); err != nil {
		return err
	}

	// Anthropic keys typically start with "sk-ant-"
	if !strings.HasPrefix(key, "sk-ant-") {
		return &ValidationError{
			Provider: "Anthropic",
			Reason:   "key should start with 'sk-ant-'",
		}
	}

	return nil
}

// ValidateOpenAIKey validates an OpenAI API key
func (v *Validator) ValidateOpenAIKey(key string) error {
	if err := v.validateBasic("OpenAI", key); err != nil {
		return err
	}

	// OpenAI keys typically start with "sk-" or "sk-proj-"
	if !strings.HasPrefix(key, "sk-") {
		return &ValidationError{
			Provider: "OpenAI",
			Reason:   "key should start with 'sk-'",
		}
	}

	return nil
}

// ValidateGoogleKey validates a Google API key
func (v *Validator) ValidateGoogleKey(key string) error {
	if err := v.validateBasic("Google", key); err != nil {
		return err
	}

	// Google API keys are typically alphanumeric with dashes
	// They don't have a consistent prefix, so we do basic validation
	if strings.Contains(key, " ") {
		return &ValidationError{
			Provider: "Google",
			Reason:   "key should not contain spaces",
		}
	}

	return nil
}

// validateBasic performs basic validation common to all API keys
func (v *Validator) validateBasic(provider, key string) error {
	if key == "" {
		return &ValidationError{
			Provider: provider,
			Reason:   "key is empty",
		}
	}

	if len(key) < v.minLength {
		return &ValidationError{
			Provider: provider,
			Reason:   fmt.Sprintf("key too short (minimum %d characters)", v.minLength),
		}
	}

	if len(key) > v.maxLength {
		return &ValidationError{
			Provider: provider,
			Reason:   fmt.Sprintf("key too long (maximum %d characters)", v.maxLength),
		}
	}

	// Check for common placeholder values
	placeholders := []string{
		"your-api-key",
		"your_api_key",
		"api-key-here",
		"placeholder",
		"xxx",
		"test",
		"example",
	}

	lowerKey := strings.ToLower(key)
	for _, placeholder := range placeholders {
		if lowerKey == placeholder || strings.Contains(lowerKey, placeholder) {
			return &ValidationError{
				Provider: provider,
				Reason:   "key appears to be a placeholder value",
			}
		}
	}

	return nil
}

// MaskKey masks an API key for safe logging
// Shows first 2 and last 2 characters for better security
func MaskKey(key string) string {
	if key == "" {
		return "<empty>"
	}

	length := len(key)
	if length <= 4 {
		return "****"
	}

	return fmt.Sprintf("%s...%s", key[:2], key[length-2:])
}

// ValidateAll validates all provided credentials and returns any errors
func (v *Validator) ValidateAll(anthropicKey, openaiKey, googleKey string) error {
	var errs []error

	if anthropicKey != "" {
		if err := v.ValidateAnthropicKey(anthropicKey); err != nil {
			errs = append(errs, err)
		}
	}

	if openaiKey != "" {
		if err := v.ValidateOpenAIKey(openaiKey); err != nil {
			errs = append(errs, err)
		}
	}

	if googleKey != "" {
		if err := v.ValidateGoogleKey(googleKey); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

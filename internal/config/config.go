package config

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Config holds all server configuration
type Config struct {
	// Provider API keys
	AnthropicAPIKey string
	OpenAIAPIKey    string
	GoogleAPIKey    string

	// Server settings
	LogLevel        string
	DefaultProvider string
	ReviewTimeout   time.Duration
	MaxDiffSize     int

	// Per-provider timeouts
	AnthropicTimeout time.Duration
	OpenAITimeout    time.Duration
	GoogleTimeout    time.Duration
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		AnthropicAPIKey:  os.Getenv("ANTHROPIC_API_KEY"),
		OpenAIAPIKey:     os.Getenv("OPENAI_API_KEY"),
		GoogleAPIKey:     os.Getenv("GOOGLE_API_KEY"),
		LogLevel:         GetEnvWithFallback("MCP_PR_LOG_LEVEL", "MCP_LOG_LEVEL", "info"),
		DefaultProvider:  GetEnvWithFallback("MCP_PR_DEFAULT_PROVIDER", "MCP_DEFAULT_PROVIDER", "anthropic"),
		ReviewTimeout:    parseDuration(GetEnvWithFallback("MCP_PR_REVIEW_TIMEOUT", "MCP_REVIEW_TIMEOUT", "120s"), 120*time.Second),
		MaxDiffSize:      parseInt(GetEnvWithFallback("MCP_PR_MAX_DIFF_SIZE", "MCP_MAX_DIFF_SIZE", "10000"), 10000),
		AnthropicTimeout: parseDuration(getEnv("ANTHROPIC_TIMEOUT", "90s"), 90*time.Second),
		OpenAITimeout:    parseDuration(getEnv("OPENAI_TIMEOUT", "90s"), 90*time.Second),
		GoogleTimeout:    parseDuration(getEnv("GOOGLE_TIMEOUT", "90s"), 90*time.Second),
	}

	// Validate at least one API key is present
	if cfg.AnthropicAPIKey == "" && cfg.OpenAIAPIKey == "" && cfg.GoogleAPIKey == "" {
		return nil, fmt.Errorf("at least one provider API key must be configured (ANTHROPIC_API_KEY, OPENAI_API_KEY, or GOOGLE_API_KEY)")
	}

	return cfg, nil
}

// HasProvider checks if a provider is configured
func (c *Config) HasProvider(provider string) bool {
	switch provider {
	case "anthropic":
		return c.AnthropicAPIKey != ""
	case "openai":
		return c.OpenAIAPIKey != ""
	case "google":
		return c.GoogleAPIKey != ""
	default:
		return false
	}
}

// GetEnvWithFallback gets environment variable with backward compatibility fallback
// It checks the new variable name first, then falls back to the old name with a deprecation warning,
// and finally uses the default value if neither is set.
func GetEnvWithFallback(newKey, oldKey, defaultValue string) string {
	// 1. Check new variable name first
	if value := os.Getenv(newKey); value != "" {
		return value
	}

	// 2. Fall back to old variable name with deprecation warning
	if value := os.Getenv(oldKey); value != "" {
		log.Printf("WARN: Environment variable %q is deprecated and will be removed in v1.0.0. Please use %q instead.", oldKey, newKey)
		return value
	}

	// 3. Use default value
	return defaultValue
}

// getEnv gets environment variable with default fallback
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseDuration parses duration string with fallback
func parseDuration(value string, defaultValue time.Duration) time.Duration {
	if d, err := time.ParseDuration(value); err == nil {
		return d
	}
	return defaultValue
}

// parseInt parses int string with fallback
func parseInt(value string, defaultValue int) int {
	var result int
	if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
		return result
	}
	return defaultValue
}

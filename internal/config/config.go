package config

import (
	"fmt"
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
		LogLevel:         getEnv("MCP_LOG_LEVEL", "info"),
		DefaultProvider:  getEnv("MCP_DEFAULT_PROVIDER", "anthropic"),
		ReviewTimeout:    parseDuration(getEnv("MCP_REVIEW_TIMEOUT", "120s"), 120*time.Second),
		MaxDiffSize:      parseInt(getEnv("MCP_MAX_DIFF_SIZE", "10000"), 10000),
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

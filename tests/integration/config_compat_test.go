package integration_test

import (
	"os"
	"testing"

	"github.com/dshills/mcp-pr/internal/config"
)

// TestConfigBackwardCompatibility tests that old environment variables work with deprecation warnings (T011)
func TestConfigBackwardCompatibility(t *testing.T) {
	t.Run("old variables work with warnings", func(t *testing.T) {
		// Clean environment
		os.Clearenv()

		// Set only old variable names
		os.Setenv("MCP_LOG_LEVEL", "warn")
		os.Setenv("MCP_DEFAULT_PROVIDER", "google")
		os.Setenv("MCP_REVIEW_TIMEOUT", "90s")
		os.Setenv("MCP_MAX_DIFF_SIZE", "5000")
		os.Setenv("ANTHROPIC_API_KEY", "test-key")

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			t.Fatalf("config.Load() failed: %v", err)
		}

		// Verify old variables are used
		if cfg.LogLevel != "warn" {
			t.Errorf("Expected LogLevel 'warn', got %q", cfg.LogLevel)
		}
		if cfg.DefaultProvider != "google" {
			t.Errorf("Expected DefaultProvider 'google', got %q", cfg.DefaultProvider)
		}
	})

	t.Run("new variables work without warnings", func(t *testing.T) {
		// Clean environment
		os.Clearenv()

		// Set only new variable names
		os.Setenv("MCP_PR_LOG_LEVEL", "debug")
		os.Setenv("MCP_PR_DEFAULT_PROVIDER", "openai")
		os.Setenv("MCP_PR_REVIEW_TIMEOUT", "180s")
		os.Setenv("MCP_PR_MAX_DIFF_SIZE", "20000")
		os.Setenv("ANTHROPIC_API_KEY", "test-key")

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			t.Fatalf("config.Load() failed: %v", err)
		}

		// Verify new variables are used
		if cfg.LogLevel != "debug" {
			t.Errorf("Expected LogLevel 'debug', got %q", cfg.LogLevel)
		}
		if cfg.DefaultProvider != "openai" {
			t.Errorf("Expected DefaultProvider 'openai', got %q", cfg.DefaultProvider)
		}
	})

	t.Run("mixed old/new configuration", func(t *testing.T) {
		// Clean environment
		os.Clearenv()

		// Mix old and new variable names
		os.Setenv("MCP_PR_LOG_LEVEL", "debug")      // New
		os.Setenv("MCP_DEFAULT_PROVIDER", "google") // Old
		os.Setenv("MCP_PR_REVIEW_TIMEOUT", "180s")  // New
		os.Setenv("MCP_MAX_DIFF_SIZE", "5000")      // Old
		os.Setenv("ANTHROPIC_API_KEY", "test-key")

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			t.Fatalf("config.Load() failed: %v", err)
		}

		// Verify correct precedence
		if cfg.LogLevel != "debug" {
			t.Errorf("Expected LogLevel 'debug' (new var), got %q", cfg.LogLevel)
		}
		if cfg.DefaultProvider != "google" {
			t.Errorf("Expected DefaultProvider 'google' (old var), got %q", cfg.DefaultProvider)
		}
	})
}

package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/dshills/mcp-pr/internal/config"
)

// TestGetEnvWithFallback tests the backward-compatible environment variable loading
func TestGetEnvWithFallback(t *testing.T) {
	tests := []struct {
		name          string
		newKey        string
		oldKey        string
		defaultValue  string
		envVars       map[string]string
		expectedValue string
		expectWarning bool
	}{
		{
			name:          "new variable only",
			newKey:        "TEST_NEW_VAR",
			oldKey:        "TEST_OLD_VAR",
			defaultValue:  "default",
			envVars:       map[string]string{"TEST_NEW_VAR": "new_value"},
			expectedValue: "new_value",
			expectWarning: false,
		},
		{
			name:          "old variable only (backward compat)",
			newKey:        "TEST_NEW_VAR",
			oldKey:        "TEST_OLD_VAR",
			defaultValue:  "default",
			envVars:       map[string]string{"TEST_OLD_VAR": "old_value"},
			expectedValue: "old_value",
			expectWarning: true,
		},
		{
			name:          "both set - new takes precedence",
			newKey:        "TEST_NEW_VAR",
			oldKey:        "TEST_OLD_VAR",
			defaultValue:  "default",
			envVars:       map[string]string{"TEST_NEW_VAR": "new_value", "TEST_OLD_VAR": "old_value"},
			expectedValue: "new_value",
			expectWarning: false,
		},
		{
			name:          "neither set - use default",
			newKey:        "TEST_NEW_VAR",
			oldKey:        "TEST_OLD_VAR",
			defaultValue:  "default",
			envVars:       map[string]string{},
			expectedValue: "default",
			expectWarning: false,
		},
		{
			name:          "empty string values - old var",
			newKey:        "TEST_NEW_VAR",
			oldKey:        "TEST_OLD_VAR",
			defaultValue:  "default",
			envVars:       map[string]string{"TEST_OLD_VAR": ""},
			expectedValue: "default", // Empty string should fall through to default
			expectWarning: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean environment
			os.Unsetenv(tt.newKey)
			os.Unsetenv(tt.oldKey)

			// Set test environment variables
			for key, value := range tt.envVars {
				if err := os.Setenv(key, value); err != nil {
					t.Fatalf("Failed to set env var %s: %v", key, err)
				}
			}

			// Clean up after test
			defer func() {
				os.Unsetenv(tt.newKey)
				os.Unsetenv(tt.oldKey)
			}()

			// Call the function (this will fail initially - function doesn't exist yet)
			result := config.GetEnvWithFallback(tt.newKey, tt.oldKey, tt.defaultValue)

			// Verify result
			if result != tt.expectedValue {
				t.Errorf("Expected %q, got %q", tt.expectedValue, result)
			}

			// Note: Deprecation warning testing will be added after we implement logging
			// For now, we're just testing the value logic
		})
	}
}

// TestConfigLoad_LogLevel tests LogLevel configuration with new and old variables (T007)
func TestConfigLoad_LogLevel(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected string
	}{
		{
			name:     "new variable MCP_PR_LOG_LEVEL",
			envVars:  map[string]string{"MCP_PR_LOG_LEVEL": "debug", "ANTHROPIC_API_KEY": "test-key"},
			expected: "debug",
		},
		{
			name:     "old variable MCP_LOG_LEVEL (backward compat)",
			envVars:  map[string]string{"MCP_LOG_LEVEL": "warn", "ANTHROPIC_API_KEY": "test-key"},
			expected: "warn",
		},
		{
			name:     "both set - new takes precedence",
			envVars:  map[string]string{"MCP_PR_LOG_LEVEL": "debug", "MCP_LOG_LEVEL": "error", "ANTHROPIC_API_KEY": "test-key"},
			expected: "debug",
		},
		{
			name:     "neither set - default value",
			envVars:  map[string]string{"ANTHROPIC_API_KEY": "test-key"},
			expected: "info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean environment
			os.Clearenv()

			// Set test environment variables
			for key, value := range tt.envVars {
				if err := os.Setenv(key, value); err != nil {
					t.Fatalf("Failed to set env var %s: %v", key, err)
				}
			}

			// Load configuration
			cfg, err := config.Load()
			if err != nil {
				t.Fatalf("config.Load() failed: %v", err)
			}

			// Verify LogLevel
			if cfg.LogLevel != tt.expected {
				t.Errorf("Expected LogLevel %q, got %q", tt.expected, cfg.LogLevel)
			}
		})
	}
}

// TestConfigLoad_DefaultProvider tests DefaultProvider configuration (T008)
func TestConfigLoad_DefaultProvider(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected string
	}{
		{
			name:     "new variable MCP_PR_DEFAULT_PROVIDER",
			envVars:  map[string]string{"MCP_PR_DEFAULT_PROVIDER": "openai", "ANTHROPIC_API_KEY": "test-key"},
			expected: "openai",
		},
		{
			name:     "old variable MCP_DEFAULT_PROVIDER (backward compat)",
			envVars:  map[string]string{"MCP_DEFAULT_PROVIDER": "google", "ANTHROPIC_API_KEY": "test-key"},
			expected: "google",
		},
		{
			name:     "both set - new takes precedence",
			envVars:  map[string]string{"MCP_PR_DEFAULT_PROVIDER": "openai", "MCP_DEFAULT_PROVIDER": "google", "ANTHROPIC_API_KEY": "test-key"},
			expected: "openai",
		},
		{
			name:     "neither set - default value",
			envVars:  map[string]string{"ANTHROPIC_API_KEY": "test-key"},
			expected: "anthropic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean environment
			os.Clearenv()

			// Set test environment variables
			for key, value := range tt.envVars {
				if err := os.Setenv(key, value); err != nil {
					t.Fatalf("Failed to set env var %s: %v", key, err)
				}
			}

			// Load configuration
			cfg, err := config.Load()
			if err != nil {
				t.Fatalf("config.Load() failed: %v", err)
			}

			// Verify DefaultProvider
			if cfg.DefaultProvider != tt.expected {
				t.Errorf("Expected DefaultProvider %q, got %q", tt.expected, cfg.DefaultProvider)
			}
		})
	}
}

// TestConfigLoad_ReviewTimeout tests ReviewTimeout configuration (T009)
func TestConfigLoad_ReviewTimeout(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected time.Duration
	}{
		{
			name:     "new variable MCP_PR_REVIEW_TIMEOUT",
			envVars:  map[string]string{"MCP_PR_REVIEW_TIMEOUT": "180s", "ANTHROPIC_API_KEY": "test-key"},
			expected: 180 * time.Second,
		},
		{
			name:     "old variable MCP_REVIEW_TIMEOUT (backward compat)",
			envVars:  map[string]string{"MCP_REVIEW_TIMEOUT": "90s", "ANTHROPIC_API_KEY": "test-key"},
			expected: 90 * time.Second,
		},
		{
			name:     "both set - new takes precedence",
			envVars:  map[string]string{"MCP_PR_REVIEW_TIMEOUT": "180s", "MCP_REVIEW_TIMEOUT": "90s", "ANTHROPIC_API_KEY": "test-key"},
			expected: 180 * time.Second,
		},
		{
			name:     "neither set - default value",
			envVars:  map[string]string{"ANTHROPIC_API_KEY": "test-key"},
			expected: 120 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean environment
			os.Clearenv()

			// Set test environment variables
			for key, value := range tt.envVars {
				if err := os.Setenv(key, value); err != nil {
					t.Fatalf("Failed to set env var %s: %v", key, err)
				}
			}

			// Load configuration
			cfg, err := config.Load()
			if err != nil {
				t.Fatalf("config.Load() failed: %v", err)
			}

			// Verify ReviewTimeout
			if cfg.ReviewTimeout != tt.expected {
				t.Errorf("Expected ReviewTimeout %v, got %v", tt.expected, cfg.ReviewTimeout)
			}
		})
	}
}

// TestConfigLoad_MaxDiffSize tests MaxDiffSize configuration (T010)
func TestConfigLoad_MaxDiffSize(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected int
	}{
		{
			name:     "new variable MCP_PR_MAX_DIFF_SIZE",
			envVars:  map[string]string{"MCP_PR_MAX_DIFF_SIZE": "20000", "ANTHROPIC_API_KEY": "test-key"},
			expected: 20000,
		},
		{
			name:     "old variable MCP_MAX_DIFF_SIZE (backward compat)",
			envVars:  map[string]string{"MCP_MAX_DIFF_SIZE": "5000", "ANTHROPIC_API_KEY": "test-key"},
			expected: 5000,
		},
		{
			name:     "both set - new takes precedence",
			envVars:  map[string]string{"MCP_PR_MAX_DIFF_SIZE": "20000", "MCP_MAX_DIFF_SIZE": "5000", "ANTHROPIC_API_KEY": "test-key"},
			expected: 20000,
		},
		{
			name:     "neither set - default value",
			envVars:  map[string]string{"ANTHROPIC_API_KEY": "test-key"},
			expected: 10000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean environment
			os.Clearenv()

			// Set test environment variables
			for key, value := range tt.envVars {
				if err := os.Setenv(key, value); err != nil {
					t.Fatalf("Failed to set env var %s: %v", key, err)
				}
			}

			// Load configuration
			cfg, err := config.Load()
			if err != nil {
				t.Fatalf("config.Load() failed: %v", err)
			}

			// Verify MaxDiffSize
			if cfg.MaxDiffSize != tt.expected {
				t.Errorf("Expected MaxDiffSize %d, got %d", tt.expected, cfg.MaxDiffSize)
			}
		})
	}
}

// TestConfigLoad_APIKeys tests that API keys load with standard names unchanged (T019)
func TestConfigLoad_APIKeys(t *testing.T) {
	tests := []struct {
		name          string
		envVars       map[string]string
		expectError   bool
		checkAntropic bool
		checkOpenAI   bool
		checkGoogle   bool
	}{
		{
			name: "ANTHROPIC_API_KEY loads unchanged",
			envVars: map[string]string{
				"ANTHROPIC_API_KEY": "sk-ant-test-key",
			},
			expectError:   false,
			checkAntropic: true,
		},
		{
			name: "OPENAI_API_KEY loads unchanged",
			envVars: map[string]string{
				"OPENAI_API_KEY": "sk-openai-test-key",
			},
			expectError: false,
			checkOpenAI: true,
		},
		{
			name: "GOOGLE_API_KEY loads unchanged",
			envVars: map[string]string{
				"GOOGLE_API_KEY": "google-test-key",
			},
			expectError: false,
			checkGoogle: true,
		},
		{
			name:        "no API keys - should error",
			envVars:     map[string]string{},
			expectError: true,
		},
		{
			name: "all three API keys work together",
			envVars: map[string]string{
				"ANTHROPIC_API_KEY": "sk-ant-test",
				"OPENAI_API_KEY":    "sk-openai-test",
				"GOOGLE_API_KEY":    "google-test",
			},
			expectError:   false,
			checkAntropic: true,
			checkOpenAI:   true,
			checkGoogle:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			for key, value := range tt.envVars {
				if err := os.Setenv(key, value); err != nil {
					t.Fatalf("Failed to set env var %s: %v", key, err)
				}
			}
			cfg, err := config.Load()
			if tt.expectError {
				if err == nil {
					t.Fatal("Expected error when no API keys set")
				}
				return
			}
			if err != nil {
				t.Fatalf("config.Load() failed: %v", err)
			}
			if tt.checkAntropic && cfg.AnthropicAPIKey == "" {
				t.Error("Expected Anthropic API key to be set")
			}
			if tt.checkOpenAI && cfg.OpenAIAPIKey == "" {
				t.Error("Expected OpenAI API key to be set")
			}
			if tt.checkGoogle && cfg.GoogleAPIKey == "" {
				t.Error("Expected Google API key to be set")
			}
		})
	}
}

// TestConfigLoad_ProviderTimeouts tests that provider timeouts load unchanged (T020)
func TestConfigLoad_ProviderTimeouts(t *testing.T) {
	tests := []struct {
		name             string
		envVars          map[string]string
		expectedAntropic time.Duration
		expectedOpenAI   time.Duration
		expectedGoogle   time.Duration
	}{
		{
			name: "ANTHROPIC_TIMEOUT loads unchanged",
			envVars: map[string]string{
				"ANTHROPIC_API_KEY": "test-key",
				"ANTHROPIC_TIMEOUT": "120s",
			},
			expectedAntropic: 120 * time.Second,
			expectedOpenAI:   90 * time.Second,
			expectedGoogle:   90 * time.Second,
		},
		{
			name: "default values when not set",
			envVars: map[string]string{
				"ANTHROPIC_API_KEY": "test-key",
			},
			expectedAntropic: 90 * time.Second,
			expectedOpenAI:   90 * time.Second,
			expectedGoogle:   90 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			for key, value := range tt.envVars {
				if err := os.Setenv(key, value); err != nil {
					t.Fatalf("Failed to set env var %s: %v", key, err)
				}
			}
			cfg, err := config.Load()
			if err != nil {
				t.Fatalf("config.Load() failed: %v", err)
			}
			if cfg.AnthropicTimeout != tt.expectedAntropic {
				t.Errorf("Expected AnthropicTimeout %v, got %v", tt.expectedAntropic, cfg.AnthropicTimeout)
			}
			if cfg.OpenAITimeout != tt.expectedOpenAI {
				t.Errorf("Expected OpenAITimeout %v, got %v", tt.expectedOpenAI, cfg.OpenAITimeout)
			}
			if cfg.GoogleTimeout != tt.expectedGoogle {
				t.Errorf("Expected GoogleTimeout %v, got %v", tt.expectedGoogle, cfg.GoogleTimeout)
			}
		})
	}
}

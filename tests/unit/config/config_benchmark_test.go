package config_test

import (
	"os"
	"testing"

	"github.com/dshills/mcp-pr/internal/config"
)

// BenchmarkConfigLoad benchmarks the configuration loading performance (T018)
func BenchmarkConfigLoad(b *testing.B) {
	// Set up test environment with minimal configuration
	_ = os.Setenv("ANTHROPIC_API_KEY", "test-key-for-benchmark")
	_ = os.Setenv("MCP_PR_LOG_LEVEL", "info")
	_ = os.Setenv("MCP_PR_DEFAULT_PROVIDER", "anthropic")
	_ = os.Setenv("MCP_PR_REVIEW_TIMEOUT", "120s")
	_ = os.Setenv("MCP_PR_MAX_DIFF_SIZE", "10000")

	defer func() {
		_ = os.Unsetenv("ANTHROPIC_API_KEY")
		_ = os.Unsetenv("MCP_PR_LOG_LEVEL")
		_ = os.Unsetenv("MCP_PR_DEFAULT_PROVIDER")
		_ = os.Unsetenv("MCP_PR_REVIEW_TIMEOUT")
		_ = os.Unsetenv("MCP_PR_MAX_DIFF_SIZE")
	}()

	// Reset timer to exclude setup time
	b.ResetTimer()

	// Run benchmark
	for i := 0; i < b.N; i++ {
		_, err := config.Load()
		if err != nil {
			b.Fatalf("config.Load() failed: %v", err)
		}
	}
}

// BenchmarkConfigLoadWithBackwardCompat benchmarks configuration loading with backward compatibility
func BenchmarkConfigLoadWithBackwardCompat(b *testing.B) {
	// Set up test environment with old variable names (worst case - triggers warnings)
	_ = os.Setenv("ANTHROPIC_API_KEY", "test-key-for-benchmark")
	_ = os.Setenv("MCP_LOG_LEVEL", "info")
	_ = os.Setenv("MCP_DEFAULT_PROVIDER", "anthropic")
	_ = os.Setenv("MCP_REVIEW_TIMEOUT", "120s")
	_ = os.Setenv("MCP_MAX_DIFF_SIZE", "10000")

	defer func() {
		_ = os.Unsetenv("ANTHROPIC_API_KEY")
		_ = os.Unsetenv("MCP_LOG_LEVEL")
		_ = os.Unsetenv("MCP_DEFAULT_PROVIDER")
		_ = os.Unsetenv("MCP_REVIEW_TIMEOUT")
		_ = os.Unsetenv("MCP_MAX_DIFF_SIZE")
	}()

	// Reset timer to exclude setup time
	b.ResetTimer()

	// Run benchmark
	for i := 0; i < b.N; i++ {
		_, err := config.Load()
		if err != nil {
			b.Fatalf("config.Load() failed: %v", err)
		}
	}
}

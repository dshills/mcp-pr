package integration

import (
	"os"
	"testing"
	"time"

	"github.com/dshills/mcp-pr/internal/providers"
)

// TestAnthropicIntegration tests real Anthropic API calls
func TestAnthropicIntegration(t *testing.T) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("ANTHROPIC_API_KEY not set, skipping integration test")
	}

	provider, err := providers.NewAnthropicProvider(apiKey, 30*time.Second)
	if err != nil {
		t.Fatalf("NewAnthropicProvider() error = %v", err)
	}
	testProviderIntegration(t, provider, "anthropic")
}

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

	provider := providers.NewAnthropicProvider(apiKey, 30*time.Second)
	testProviderIntegration(t, provider, "anthropic")
}

package integration

import (
	"os"
	"testing"
	"time"

	"github.com/dshills/mcp-pr/internal/providers"
)

// TestOpenAIIntegration tests real OpenAI API calls
func TestOpenAIIntegration(t *testing.T) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Skip("OPENAI_API_KEY not set, skipping integration test")
	}

	provider := providers.NewOpenAIProvider(apiKey, 30*time.Second)
	testProviderIntegration(t, provider, "openai")
}

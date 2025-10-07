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

	provider, err := providers.NewOpenAIProvider(apiKey, 30*time.Second)
	if err != nil {
		t.Fatalf("NewOpenAIProvider() error = %v", err)
	}
	testProviderIntegration(t, provider, "openai")
}

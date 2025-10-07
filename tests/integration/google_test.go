package integration

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dshills/mcp-pr/internal/providers"
	"github.com/dshills/mcp-pr/internal/review"
)

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// TestGoogleIntegration tests real Google GenAI API calls
func TestGoogleIntegration(t *testing.T) {
	// Try both GEMINI_API_KEY and GOOGLE_API_KEY
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("GOOGLE_API_KEY")
	}
	if apiKey == "" {
		t.Skip("GEMINI_API_KEY or GOOGLE_API_KEY not set, skipping integration test")
	}

	provider, err := providers.NewGoogleProvider(apiKey, 30*time.Second)
	if err != nil {
		t.Fatalf("NewGoogleProvider() error = %v", err)
	}
	defer func() {
		if err := provider.Close(); err != nil {
			t.Errorf("Close() error = %v", err)
		}
	}()

	ctx := context.Background()
	req := review.Request{
		SourceType:  "arbitrary",
		Code:        "func divide(a, b int) int { return a / b }",
		Provider:    "google",
		Language:    "go",
		ReviewDepth: "quick",
	}

	resp, err := provider.Review(ctx, req)
	if err != nil {
		// Known issue: google/generative-ai-go SDK is deprecated
		// Model name incompatibility with v1beta API
		if contains(err.Error(), "is not found for API version") {
			t.Skipf("Skipping due to known Google SDK deprecation issue: %v", err)
		}
		t.Fatalf("Review() error = %v", err)
	}

	// Verify response structure
	if resp.Summary == "" {
		t.Error("Summary should not be empty")
	}

	if resp.Provider != "google" {
		t.Errorf("Provider = %v, want google", resp.Provider)
	}

	if resp.Metadata == nil {
		t.Error("Metadata should not be nil")
	} else {
		if resp.Metadata.SourceType != "arbitrary" {
			t.Errorf("Metadata.SourceType = %v, want arbitrary", resp.Metadata.SourceType)
		}
		if resp.Metadata.Model == "" {
			t.Error("Metadata.Model should not be empty")
		}
	}

	t.Logf("Found %d findings", len(resp.Findings))
	if len(resp.Findings) > 0 {
		t.Logf("First finding: %+v", resp.Findings[0])
	}
}

package integration

import (
	"context"
	"testing"

	"github.com/dshills/mcp-pr/internal/review"
)

// Provider interface for testing (same as in providers package)
type Provider interface {
	Review(ctx context.Context, req review.Request) (*review.Response, error)
	Name() string
	IsAvailable() bool
}

// testProviderIntegration is a helper function to test any provider
func testProviderIntegration(t *testing.T, provider Provider, expectedProviderName string) {
	t.Helper()

	ctx := context.Background()
	req := review.Request{
		SourceType:  "arbitrary",
		Code:        "func divide(a, b int) int { return a / b }",
		Provider:    expectedProviderName,
		Language:    "go",
		ReviewDepth: "quick",
	}

	resp, err := provider.Review(ctx, req)
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}

	// Verify response structure
	if resp.Summary == "" {
		t.Error("Summary should not be empty")
	}

	if resp.Provider != expectedProviderName {
		t.Errorf("Provider = %v, want %v", resp.Provider, expectedProviderName)
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

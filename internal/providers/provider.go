package providers

import (
	"context"
	"github.com/dshills/mcp-pr/internal/review"
)

// Provider defines the interface for LLM providers
type Provider interface {
	// Review analyzes code and returns structured findings
	Review(ctx context.Context, req review.Request) (*review.Response, error)

	// Name returns the provider name ("anthropic", "openai", "google")
	Name() string

	// IsAvailable checks if the provider is configured and ready
	IsAvailable() bool
}

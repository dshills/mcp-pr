package review

import (
	"context"
	"fmt"
	"time"

	"github.com/dshills/mcp-pr/internal/logging"
)

// Provider interface for LLM providers (avoid circular import)
type Provider interface {
	Review(ctx context.Context, req Request) (*Response, error)
	Name() string
	IsAvailable() bool
}

// Engine orchestrates code review operations
type Engine struct {
	providers       map[string]Provider
	defaultProvider string
	maxRetries      int
	retryDelay      time.Duration
}

// NewEngine creates a new review engine
func NewEngine(providers map[string]Provider, defaultProvider string) *Engine {
	return &Engine{
		providers:       providers,
		defaultProvider: defaultProvider,
		maxRetries:      3,
		retryDelay:      time.Second,
	}
}

// Review performs a code review using the specified or default provider
func (e *Engine) Review(ctx context.Context, req Request) (*Response, error) {
	start := time.Now()

	// Validate request
	if err := req.Validate(); err != nil {
		logging.Error(ctx, "Invalid review request", "error", err)
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Select provider
	providerName := req.Provider
	if providerName == "" {
		providerName = e.defaultProvider
	}

	provider, exists := e.providers[providerName]
	if !exists {
		logging.Error(ctx, "Provider not found", "provider", providerName)
		return nil, fmt.Errorf("provider %s not found", providerName)
	}

	if !provider.IsAvailable() {
		logging.Error(ctx, "Provider not available", "provider", providerName)
		return nil, fmt.Errorf("provider %s not available", providerName)
	}

	logging.Info(ctx, "Starting code review",
		"provider", providerName,
		"source_type", req.SourceType,
		"language", req.Language,
		"review_depth", req.ReviewDepth,
	)

	// Perform review with retry logic
	var resp *Response
	var err error

	for attempt := 0; attempt <= e.maxRetries; attempt++ {
		if attempt > 0 {
			logging.Info(ctx, "Retrying review",
				"attempt", attempt,
				"provider", providerName,
			)
			time.Sleep(e.retryDelay * time.Duration(attempt))
		}

		resp, err = provider.Review(ctx, req)
		if err == nil {
			break
		}

		logging.Warn(ctx, "Review attempt failed",
			"attempt", attempt,
			"provider", providerName,
			"error", err,
		)
	}

	if err != nil {
		logging.Error(ctx, "Review failed after retries",
			"provider", providerName,
			"attempts", e.maxRetries+1,
			"error", err,
		)
		return nil, fmt.Errorf("review failed after %d attempts: %w", e.maxRetries+1, err)
	}

	duration := time.Since(start)
	logging.Info(ctx, "Review completed",
		"provider", providerName,
		"findings_count", len(resp.Findings),
		"duration_ms", duration.Milliseconds(),
	)

	return resp, nil
}

// GetProvider returns a provider by name
func (e *Engine) GetProvider(name string) (Provider, bool) {
	provider, exists := e.providers[name]
	return provider, exists
}

// ListProviders returns names of all available providers
func (e *Engine) ListProviders() []string {
	names := make([]string, 0, len(e.providers))
	for name, provider := range e.providers {
		if provider.IsAvailable() {
			names = append(names, name)
		}
	}
	return names
}

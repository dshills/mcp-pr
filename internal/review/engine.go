package review

import (
	"context"
	"fmt"
	"time"

	"github.com/dshills/mcp-pr/internal/git"
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
	maxDiffSize     int
}

// NewEngine creates a new review engine
func NewEngine(providers map[string]Provider, defaultProvider string, maxDiffSize int) *Engine {
	return &Engine{
		providers:       providers,
		defaultProvider: defaultProvider,
		maxRetries:      1, // Reduced from 3 to 1 to avoid long delays
		retryDelay:      time.Second,
		maxDiffSize:     maxDiffSize,
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

	// Populate Code field from git if needed
	if err := e.populateCodeFromGit(ctx, &req); err != nil {
		logging.Error(ctx, "Failed to get git diff", "error", err)
		return nil, fmt.Errorf("failed to get git diff: %w", err)
	}

	// Validate diff size
	if len(req.Code) > e.maxDiffSize {
		logging.Error(ctx, "Diff too large",
			"size_bytes", len(req.Code),
			"max_size_bytes", e.maxDiffSize,
		)
		return nil, fmt.Errorf("diff size (%d bytes) exceeds maximum allowed size (%d bytes). Consider reviewing smaller changes or increasing MCP_MAX_DIFF_SIZE",
			len(req.Code), e.maxDiffSize)
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
		"code_size_bytes", len(req.Code),
	)

	// Perform review with retry logic
	var resp *Response
	var err error

	for attempt := 0; attempt <= e.maxRetries; attempt++ {
		if attempt > 0 {
			logging.Info(ctx, "Retrying review",
				"attempt", attempt,
				"max_retries", e.maxRetries,
				"provider", providerName,
			)
			time.Sleep(e.retryDelay * time.Duration(attempt))
		}

		logging.Info(ctx, "Sending review request to LLM",
			"provider", providerName,
			"attempt", attempt+1,
			"max_attempts", e.maxRetries+1,
		)

		resp, err = provider.Review(ctx, req)
		if err == nil {
			logging.Info(ctx, "Review request completed successfully",
				"provider", providerName,
				"attempt", attempt+1,
			)
			break
		}

		logging.Warn(ctx, "Review attempt failed",
			"attempt", attempt+1,
			"max_attempts", e.maxRetries+1,
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

// populateCodeFromGit retrieves git diff and populates the Code field
func (e *Engine) populateCodeFromGit(ctx context.Context, req *Request) error {
	// Skip if not a git-based request
	if req.SourceType == "arbitrary" {
		return nil
	}

	// Skip if Code is already populated
	if req.Code != "" {
		return nil
	}

	logging.Info(ctx, "Fetching git diff",
		"source_type", req.SourceType,
		"repository", req.RepositoryPath,
	)

	// Create git client
	client := git.NewClient(req.RepositoryPath)

	// Get appropriate diff based on source type
	var diff string
	var err error

	switch req.SourceType {
	case "staged":
		diff, err = client.GetStagedDiffContext(ctx)
	case "unstaged":
		diff, err = client.GetUnstagedDiffContext(ctx)
	case "commit":
		diff, err = client.GetCommitDiffContext(ctx, req.CommitSHA)
	default:
		return fmt.Errorf("unsupported source type: %s", req.SourceType)
	}

	if err != nil {
		return err
	}

	logging.Info(ctx, "Git diff fetched",
		"diff_size_bytes", len(diff),
	)

	// Populate the Code field with diff
	req.Code = diff
	return nil
}

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dshills/mcp-pr/internal/config"
	"github.com/dshills/mcp-pr/internal/credentials"
	"github.com/dshills/mcp-pr/internal/logging"
	"github.com/dshills/mcp-pr/internal/mcp"
	"github.com/dshills/mcp-pr/internal/providers"
	"github.com/dshills/mcp-pr/internal/review"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logging
	logging.Init(cfg.LogLevel)

	ctx := context.Background()

	// Validate credentials before any logging
	validator := credentials.NewValidator()
	if err := validator.ValidateAll(cfg.AnthropicAPIKey, cfg.OpenAIAPIKey, cfg.GoogleAPIKey); err != nil {
		logging.Error(ctx, "Invalid API credentials", "error", err)
		fmt.Fprintf(os.Stderr, "Error: Invalid API credentials:\n%v\n", err)
		os.Exit(1)
	}

	logging.Info(ctx, "Starting MCP Code Review Server",
		"version", "1.0.0",
		"default_provider", cfg.DefaultProvider,
	)

	// Initialize providers
	providerMap := make(map[string]review.Provider)

	if cfg.AnthropicAPIKey != "" {
		anthropicProvider, err := providers.NewAnthropicProvider(cfg.AnthropicAPIKey, cfg.AnthropicTimeout)
		if err != nil {
			logging.Error(ctx, "Failed to initialize Anthropic provider", "error", err)
		} else {
			providerMap["anthropic"] = anthropicProvider
			logging.Info(ctx, "Initialized Anthropic provider")
		}
	}

	if cfg.OpenAIAPIKey != "" {
		openaiProvider, err := providers.NewOpenAIProvider(cfg.OpenAIAPIKey, cfg.OpenAITimeout)
		if err != nil {
			logging.Error(ctx, "Failed to initialize OpenAI provider", "error", err)
		} else {
			providerMap["openai"] = openaiProvider
			logging.Info(ctx, "Initialized OpenAI provider")
		}
	}

	if cfg.GoogleAPIKey != "" {
		googleProvider, err := providers.NewGoogleProvider(cfg.GoogleAPIKey, cfg.GoogleTimeout)
		if err != nil {
			logging.Error(ctx, "Failed to initialize Google provider", "error", err)
		} else {
			providerMap["google"] = googleProvider
			logging.Info(ctx, "Initialized Google provider")
		}
	}

	// Validate at least one provider is available
	if len(providerMap) == 0 {
		logging.Error(ctx, "No providers available - check API key configuration")
		fmt.Fprintf(os.Stderr, "Error: No LLM providers configured. Set at least one API key:\n")
		fmt.Fprintf(os.Stderr, "  ANTHROPIC_API_KEY\n")
		fmt.Fprintf(os.Stderr, "  OPENAI_API_KEY\n")
		fmt.Fprintf(os.Stderr, "  GOOGLE_API_KEY\n")
		os.Exit(1)
	}

	// Create review engine
	engine := review.NewEngine(providerMap, cfg.DefaultProvider, cfg.MaxDiffSize)
	logging.Info(ctx, "Review engine initialized",
		"providers", engine.ListProviders(),
		"max_diff_size", cfg.MaxDiffSize,
	)

	// Create MCP server
	server, err := mcp.NewServer(engine)
	if err != nil {
		logging.Error(ctx, "Failed to create MCP server", "error", err)
		fmt.Fprintf(os.Stderr, "Failed to create MCP server: %v\n", err)
		os.Exit(1)
	}

	// Run the server on stdio
	logging.Info(ctx, "Starting MCP server on stdio")
	if err := server.Run(ctx); err != nil {
		logging.Error(ctx, "Server error", "error", err)
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}

	logging.Info(ctx, "MCP Code Review Server shutting down")
}

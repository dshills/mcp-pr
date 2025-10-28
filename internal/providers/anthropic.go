package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/dshills/mcp-pr/internal/review"
)

// AnthropicProvider implements Provider for Anthropic Claude
type AnthropicProvider struct {
	client  *anthropic.Client
	timeout time.Duration
}

// NewAnthropicProvider creates a new Anthropic provider
func NewAnthropicProvider(apiKey string, timeout time.Duration) (*AnthropicProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("anthropic API key is required")
	}

	client := anthropic.NewClient(option.WithAPIKey(apiKey))
	return &AnthropicProvider{
		client:  &client,
		timeout: timeout,
	}, nil
}

// Review analyzes code using Claude
func (p *AnthropicProvider) Review(ctx context.Context, req review.Request) (*review.Response, error) {
	start := time.Now()

	// Build prompt
	prompt := buildReviewPrompt(req)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	// Call Claude API
	message, err := p.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_5, // Claude Sonnet 4.5
		MaxTokens: 4096,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("anthropic API call timed out after %v", p.timeout)
		}
		return nil, fmt.Errorf("anthropic API error: %w", err)
	}

	// Parse response
	var responseText string
	for _, block := range message.Content {
		if block.Type == "text" {
			responseText += block.Text
		}
	}

	// Parse JSON response
	findings, summary := parseReviewResponse(responseText)

	duration := time.Since(start)

	return &review.Response{
		Findings: findings,
		Summary:  summary,
		Provider: "anthropic",
		Duration: duration,
		Metadata: &review.Metadata{
			SourceType: req.SourceType,
			Model:      "claude-sonnet-4-5",
		},
	}, nil
}

// Name returns provider name
func (p *AnthropicProvider) Name() string {
	return "anthropic"
}

// IsAvailable checks if provider is configured
func (p *AnthropicProvider) IsAvailable() bool {
	return p.client != nil
}

// buildReviewPrompt creates the review prompt
func buildReviewPrompt(req review.Request) string {
	prompt := `You are a code review assistant. Analyze the following code and identify issues.

Respond in JSON format with an array of findings:
{
  "findings": [
    {
      "category": "bug|security|performance|style|best-practice",
      "severity": "critical|high|medium|low|info",
      "line": <line_number_or_null>,
      "description": "What the issue is",
      "suggestion": "How to fix it"
    }
  ],
  "summary": "Overall assessment"
}

Code to review:
` + "```" + req.Language + "\n" + req.Code + "\n```" + `

Review depth: ` + req.ReviewDepth
	return prompt
}

// parseReviewResponse extracts findings from LLM response
func parseReviewResponse(responseText string) ([]review.Finding, string) {
	type Response struct {
		Findings []review.Finding `json:"findings"`
		Summary  string           `json:"summary"`
	}

	var resp Response
	if err := json.Unmarshal([]byte(responseText), &resp); err != nil {
		// Fallback if JSON parsing fails
		return []review.Finding{}, responseText
	}

	return resp.Findings, resp.Summary
}

package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dshills/mcp-pr/internal/review"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GoogleProvider implements Provider for Google Gemini
type GoogleProvider struct {
	client  *genai.Client
	timeout time.Duration
}

// NewGoogleProvider creates a new Google provider
// Note: The google/generative-ai-go SDK is deprecated. This implementation
// may require migration to github.com/googleapis/go-genai for latest models.
// Current known issue: Model names may not work with v1beta API version.
func NewGoogleProvider(apiKey string, timeout time.Duration) (*GoogleProvider, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create google client: %w", err)
	}
	return &GoogleProvider{
		client:  client,
		timeout: timeout,
	}, nil
}

// Review analyzes code using Gemini
func (p *GoogleProvider) Review(ctx context.Context, req review.Request) (*review.Response, error) {
	start := time.Now()

	// Build prompt
	prompt := buildGoogleReviewPrompt(req)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	// Get Gemini model
	model := p.client.GenerativeModel("gemini-1.5-flash")

	// Call Gemini API
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("google API error: %w", err)
	}

	// Parse response
	var responseText string
	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
		for _, part := range resp.Candidates[0].Content.Parts {
			if textPart, ok := part.(genai.Text); ok {
				responseText += string(textPart)
			}
		}
	}

	// Parse JSON response
	findings, summary := parseGoogleReviewResponse(responseText)

	duration := time.Since(start)

	return &review.Response{
		Findings: findings,
		Summary:  summary,
		Provider: "google",
		Duration: duration,
		Metadata: &review.Metadata{
			SourceType: req.SourceType,
			Model:      "gemini-1.5-flash",
		},
	}, nil
}

// Name returns provider name
func (p *GoogleProvider) Name() string {
	return "google"
}

// IsAvailable checks if provider is configured
func (p *GoogleProvider) IsAvailable() bool {
	return p.client != nil
}

// Close closes the client connection
func (p *GoogleProvider) Close() error {
	if p.client != nil {
		return p.client.Close()
	}
	return nil
}

// buildGoogleReviewPrompt creates the review prompt for Gemini
func buildGoogleReviewPrompt(req review.Request) string {
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

// parseGoogleReviewResponse extracts findings from Gemini response
func parseGoogleReviewResponse(responseText string) ([]review.Finding, string) {
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

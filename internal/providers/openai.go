package providers

import (
	"context"
	"fmt"
	"time"

	"github.com/dshills/mcp-pr/internal/review"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// OpenAIProvider implements Provider for OpenAI GPT
type OpenAIProvider struct {
	client  *openai.Client
	timeout time.Duration
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(apiKey string, timeout time.Duration) (*OpenAIProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("openai API key is required")
	}

	client := openai.NewClient(option.WithAPIKey(apiKey))
	return &OpenAIProvider{
		client:  &client,
		timeout: timeout,
	}, nil
}

// Review analyzes code using GPT
func (p *OpenAIProvider) Review(ctx context.Context, req review.Request) (*review.Response, error) {
	start := time.Now()

	// Build system and user messages
	systemPrompt := buildSystemPrompt()
	userPrompt := buildUserPrompt(req)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	// Call OpenAI API
	chatCompletion, err := p.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(userPrompt),
		},
		Model: openai.ChatModelGPT4o,
	})

	if err != nil {
		return nil, fmt.Errorf("openai API error: %w", err)
	}

	// Parse response
	var responseText string
	if len(chatCompletion.Choices) > 0 {
		responseText = chatCompletion.Choices[0].Message.Content
	}

	// Parse JSON response
	findings, summary := parseReviewResponse(responseText)

	duration := time.Since(start)

	return &review.Response{
		Findings: findings,
		Summary:  summary,
		Provider: "openai",
		Duration: duration,
		Metadata: &review.Metadata{
			SourceType: req.SourceType,
			Model:      "gpt-4o",
		},
	}, nil
}

// Name returns provider name
func (p *OpenAIProvider) Name() string {
	return "openai"
}

// IsAvailable checks if provider is configured
func (p *OpenAIProvider) IsAvailable() bool {
	return p.client != nil
}

// buildSystemPrompt creates the system message for GPT
func buildSystemPrompt() string {
	return `You are a code review assistant. Analyze code and identify issues.

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
}`
}

// buildUserPrompt creates the user message with code
func buildUserPrompt(req review.Request) string {
	prompt := `Code to review:
` + "```" + req.Language + "\n" + req.Code + "\n```" + `

Review depth: ` + req.ReviewDepth
	return prompt
}

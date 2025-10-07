package contract

import (
	"context"
	"testing"
	"time"

	"github.com/dshills/mcp-pr/internal/providers"
	"github.com/dshills/mcp-pr/internal/review"
)

// mockProvider implements the Provider interface for testing
type mockProvider struct {
	name      string
	available bool
}

func (m *mockProvider) Review(ctx context.Context, req review.Request) (*review.Response, error) {
	// Return mock response
	return &review.Response{
		Findings: []review.Finding{},
		Summary:  "Mock review",
		Provider: m.name,
		Duration: time.Second,
	}, nil
}

func (m *mockProvider) Name() string {
	return m.name
}

func (m *mockProvider) IsAvailable() bool {
	return m.available
}

// TestProviderInterfaceContract ensures all providers implement the interface correctly
func TestProviderInterfaceContract(t *testing.T) {
	// Create mock provider
	mock := &mockProvider{
		name:      "mock",
		available: true,
	}

	// Verify it implements Provider interface
	var _ providers.Provider = mock

	// Test Name() method
	if name := mock.Name(); name != "mock" {
		t.Errorf("Name() = %v, want %v", name, "mock")
	}

	// Test IsAvailable() method
	if !mock.IsAvailable() {
		t.Error("IsAvailable() = false, want true")
	}

	// Test Review() method signature
	ctx := context.Background()
	req := review.Request{
		SourceType: "arbitrary",
		Code:       "test code",
		Provider:   "mock",
	}

	resp, err := mock.Review(ctx, req)
	if err != nil {
		t.Errorf("Review() error = %v", err)
	}

	if resp == nil {
		t.Fatal("Review() returned nil response")
	}

	if resp.Provider != "mock" {
		t.Errorf("Response.Provider = %v, want %v", resp.Provider, "mock")
	}
}

// TestAnthropicProviderContract tests that Anthropic provider implements the interface
func TestAnthropicProviderContract(t *testing.T) {
	provider := providers.NewAnthropicProvider("test-key", 30*time.Second)

	// Verify it implements Provider interface
	var _ providers.Provider = provider

	if provider.Name() != "anthropic" {
		t.Errorf("Name() = %v, want anthropic", provider.Name())
	}

	if !provider.IsAvailable() {
		t.Error("IsAvailable() should be true when client is initialized")
	}
}

// TestOpenAIProviderContract tests that OpenAI provider implements the interface
func TestOpenAIProviderContract(t *testing.T) {
	provider := providers.NewOpenAIProvider("test-key", 30*time.Second)

	// Verify it implements Provider interface
	var _ providers.Provider = provider

	if provider.Name() != "openai" {
		t.Errorf("Name() = %v, want openai", provider.Name())
	}

	if !provider.IsAvailable() {
		t.Error("IsAvailable() should be true when client is initialized")
	}
}

// TestGoogleProviderContract tests that Google provider implements the interface
func TestGoogleProviderContract(t *testing.T) {
	provider, err := providers.NewGoogleProvider("test-key", 30*time.Second)
	if err != nil {
		t.Fatalf("NewGoogleProvider() error = %v", err)
	}
	defer func() {
		if err := provider.Close(); err != nil {
			t.Errorf("Close() error = %v", err)
		}
	}()

	// Verify it implements Provider interface
	var _ providers.Provider = provider

	if provider.Name() != "google" {
		t.Errorf("Name() = %v, want google", provider.Name())
	}

	if !provider.IsAvailable() {
		t.Error("IsAvailable() should be true when client is initialized")
	}
}

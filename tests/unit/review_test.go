package unit

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dshills/mcp-pr/internal/logging"
	"github.com/dshills/mcp-pr/internal/review"
)

func init() {
	// Initialize logging for tests
	logging.Init("error") // Use error level to reduce test output
}

// mockProvider implements the Provider interface for testing
type mockProvider struct {
	name       string
	available  bool
	response   *review.Response
	err        error
	callCount  int
	reviewFunc func(ctx context.Context, req review.Request) (*review.Response, error)
}

func (m *mockProvider) Review(ctx context.Context, req review.Request) (*review.Response, error) {
	m.callCount++
	if m.reviewFunc != nil {
		return m.reviewFunc(ctx, req)
	}
	if m.err != nil {
		return nil, m.err
	}
	return m.response, nil
}

func (m *mockProvider) Name() string {
	return m.name
}

func (m *mockProvider) IsAvailable() bool {
	return m.available
}

// TestEngineReviewSuccess tests successful review with default provider
func TestEngineReviewSuccess(t *testing.T) {
	mockResp := &review.Response{
		Findings: []review.Finding{
			{
				Category:    "bug",
				Severity:    "high",
				Description: "Test finding",
			},
		},
		Summary:  "Test summary",
		Provider: "mock",
		Duration: time.Second,
	}

	provider := &mockProvider{
		name:      "mock",
		available: true,
		response:  mockResp,
	}

	providers := map[string]review.Provider{
		"mock": provider,
	}

	engine := review.NewEngine(providers, "mock")

	ctx := context.Background()
	req := review.Request{
		SourceType: "arbitrary",
		Code:       "test code",
		Provider:   "mock",
	}

	resp, err := engine.Review(ctx, req)
	if err != nil {
		t.Fatalf("Review() error = %v, want nil", err)
	}

	if resp == nil {
		t.Fatal("Review() returned nil response")
	}

	if resp.Provider != "mock" {
		t.Errorf("Provider = %v, want mock", resp.Provider)
	}

	if len(resp.Findings) != 1 {
		t.Errorf("Findings count = %v, want 1", len(resp.Findings))
	}

	if provider.callCount != 1 {
		t.Errorf("Provider called %d times, want 1", provider.callCount)
	}
}

// TestEngineReviewProviderNotFound tests error when provider doesn't exist
func TestEngineReviewProviderNotFound(t *testing.T) {
	providers := map[string]review.Provider{}
	engine := review.NewEngine(providers, "default")

	ctx := context.Background()
	req := review.Request{
		SourceType: "arbitrary",
		Code:       "test code",
		Provider:   "nonexistent",
	}

	_, err := engine.Review(ctx, req)
	if err == nil {
		t.Fatal("Review() error = nil, want error for nonexistent provider")
	}

	expectedMsg := "provider nonexistent not found"
	if err.Error() != expectedMsg {
		t.Errorf("Error = %v, want %v", err, expectedMsg)
	}
}

// TestEngineReviewProviderUnavailable tests error when provider is unavailable
func TestEngineReviewProviderUnavailable(t *testing.T) {
	provider := &mockProvider{
		name:      "mock",
		available: false,
	}

	providers := map[string]review.Provider{
		"mock": provider,
	}

	engine := review.NewEngine(providers, "mock")

	ctx := context.Background()
	req := review.Request{
		SourceType: "arbitrary",
		Code:       "test code",
		Provider:   "mock",
	}

	_, err := engine.Review(ctx, req)
	if err == nil {
		t.Fatal("Review() error = nil, want error for unavailable provider")
	}

	expectedMsg := "provider mock not available"
	if err.Error() != expectedMsg {
		t.Errorf("Error = %v, want %v", err, expectedMsg)
	}
}

// TestEngineReviewRetry tests retry logic on provider errors
func TestEngineReviewRetry(t *testing.T) {
	provider := &mockProvider{
		name:      "mock",
		available: true,
		err:       errors.New("temporary error"),
	}

	providers := map[string]review.Provider{
		"mock": provider,
	}

	engine := review.NewEngine(providers, "mock")

	ctx := context.Background()
	req := review.Request{
		SourceType: "arbitrary",
		Code:       "test code",
		Provider:   "mock",
	}

	_, err := engine.Review(ctx, req)
	if err == nil {
		t.Fatal("Review() error = nil, want error after retries")
	}

	// Should have tried 1 initial attempt + 3 retries = 4 total
	expectedCalls := 4
	if provider.callCount != expectedCalls {
		t.Errorf("Provider called %d times, want %d (1 initial + 3 retries)", provider.callCount, expectedCalls)
	}
}

// TestEngineReviewRetrySuccess tests successful review after retry
func TestEngineReviewRetrySuccess(t *testing.T) {
	mockResp := &review.Response{
		Findings: []review.Finding{},
		Summary:  "Success after retry",
		Provider: "mock",
		Duration: time.Second,
	}

	callCount := 0
	provider := &mockProvider{
		name:      "mock",
		available: true,
		reviewFunc: func(ctx context.Context, req review.Request) (*review.Response, error) {
			callCount++
			if callCount < 3 {
				return nil, errors.New("temporary error")
			}
			return mockResp, nil
		},
	}

	providers := map[string]review.Provider{
		"mock": provider,
	}

	engine := review.NewEngine(providers, "mock")

	ctx := context.Background()
	req := review.Request{
		SourceType: "arbitrary",
		Code:       "test code",
		Provider:   "mock",
	}

	resp, err := engine.Review(ctx, req)
	if err != nil {
		t.Fatalf("Review() error = %v, want nil (should succeed after retries)", err)
	}

	if resp == nil {
		t.Fatal("Review() returned nil response")
	}

	if callCount != 3 {
		t.Errorf("Provider called %d times, want 3 (2 failures + 1 success)", callCount)
	}
}

// TestEngineReviewMultipleProviders tests engine with multiple providers
func TestEngineReviewMultipleProviders(t *testing.T) {
	mockResp1 := &review.Response{
		Findings: []review.Finding{},
		Summary:  "Provider 1 response",
		Provider: "provider1",
		Duration: time.Second,
	}

	mockResp2 := &review.Response{
		Findings: []review.Finding{},
		Summary:  "Provider 2 response",
		Provider: "provider2",
		Duration: time.Second,
	}

	provider1 := &mockProvider{
		name:      "provider1",
		available: true,
		response:  mockResp1,
	}

	provider2 := &mockProvider{
		name:      "provider2",
		available: true,
		response:  mockResp2,
	}

	providers := map[string]review.Provider{
		"provider1": provider1,
		"provider2": provider2,
	}

	engine := review.NewEngine(providers, "provider1")

	ctx := context.Background()

	// Test using provider1
	req1 := review.Request{
		SourceType: "arbitrary",
		Code:       "test code",
		Provider:   "provider1",
	}

	resp1, err := engine.Review(ctx, req1)
	if err != nil {
		t.Fatalf("Review() with provider1 error = %v", err)
	}

	if resp1.Provider != "provider1" {
		t.Errorf("Provider = %v, want provider1", resp1.Provider)
	}

	// Test using provider2
	req2 := review.Request{
		SourceType: "arbitrary",
		Code:       "test code",
		Provider:   "provider2",
	}

	resp2, err := engine.Review(ctx, req2)
	if err != nil {
		t.Fatalf("Review() with provider2 error = %v", err)
	}

	if resp2.Provider != "provider2" {
		t.Errorf("Provider = %v, want provider2", resp2.Provider)
	}
}

// TestEngineReviewContextCancellation tests handling of context cancellation
func TestEngineReviewContextCancellation(t *testing.T) {
	provider := &mockProvider{
		name:      "mock",
		available: true,
		reviewFunc: func(ctx context.Context, req review.Request) (*review.Response, error) {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Second):
				return &review.Response{}, nil
			}
		},
	}

	providers := map[string]review.Provider{
		"mock": provider,
	}

	engine := review.NewEngine(providers, "mock")

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	req := review.Request{
		SourceType: "arbitrary",
		Code:       "test code",
		Provider:   "mock",
	}

	_, err := engine.Review(ctx, req)
	if err == nil {
		t.Fatal("Review() error = nil, want context cancellation error")
	}

	if !errors.Is(err, context.Canceled) {
		t.Errorf("Error = %v, want context.Canceled", err)
	}
}

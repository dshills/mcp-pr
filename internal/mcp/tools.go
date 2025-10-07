package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dshills/mcp-pr/internal/logging"
	"github.com/dshills/mcp-pr/internal/review"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// handleReviewCode handles the review_code tool request
func (s *Server) handleReviewCode(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logging.Info(ctx, "Handling review_code request")

	// Parse arguments
	var args struct {
		Code        string   `json:"code"`
		Language    string   `json:"language"`
		Provider    string   `json:"provider,omitempty"`
		ReviewDepth string   `json:"review_depth,omitempty"`
		FocusAreas  []string `json:"focus_areas,omitempty"`
	}

	if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if args.ReviewDepth == "" {
		args.ReviewDepth = "quick"
	}

	// Build review request
	reviewReq := review.Request{
		SourceType:  "arbitrary",
		Code:        args.Code,
		Provider:    args.Provider,
		Language:    args.Language,
		ReviewDepth: args.ReviewDepth,
		FocusAreas:  args.FocusAreas,
	}

	// Perform review
	resp, err := s.engine.Review(ctx, reviewReq)
	if err != nil {
		logging.Error(ctx, "Review failed", "error", err)
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Review failed: %v", err)}},
		}, nil
	}

	// Format response as JSON content
	jsonData, err := json.MarshalIndent(formatReviewResponse(resp), "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to format response: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(jsonData)}},
	}, nil
}

// handleReviewStaged handles the review_staged tool request
func (s *Server) handleReviewStaged(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return s.handleGitReview(ctx, req, "staged")
}

// handleReviewUnstaged handles the review_unstaged tool request
func (s *Server) handleReviewUnstaged(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return s.handleGitReview(ctx, req, "unstaged")
}

// handleReviewCommit handles the review_commit tool request
func (s *Server) handleReviewCommit(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logging.Info(ctx, "Handling review_commit request")

	// Parse arguments
	var args struct {
		RepositoryPath string `json:"repository_path"`
		CommitSHA      string `json:"commit_sha"`
		Provider       string `json:"provider,omitempty"`
		ReviewDepth    string `json:"review_depth,omitempty"`
	}

	if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	// Validate required arguments
	if args.RepositoryPath == "" {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: "repository_path is required"}},
		}, nil
	}

	if args.CommitSHA == "" {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: "commit_sha is required"}},
		}, nil
	}

	if args.ReviewDepth == "" {
		args.ReviewDepth = "quick"
	}

	// Build review request
	reviewReq := review.Request{
		SourceType:     "commit",
		RepositoryPath: args.RepositoryPath,
		CommitSHA:      args.CommitSHA,
		Provider:       args.Provider,
		ReviewDepth:    args.ReviewDepth,
		Language:       "diff",
	}

	// Perform review (engine will populate Code from git)
	resp, err := s.engine.Review(ctx, reviewReq)
	if err != nil {
		logging.Error(ctx, "Review failed", "error", err)
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Review failed: %v", err)}},
		}, nil
	}

	// Format response as JSON content
	jsonData, err := json.MarshalIndent(formatReviewResponse(resp), "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to format response: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(jsonData)}},
	}, nil
}

// formatReviewResponse formats the review response for output
func formatReviewResponse(resp *review.Response) map[string]interface{} {
	// Convert findings to map format
	findings := make([]map[string]interface{}, len(resp.Findings))
	for i, f := range resp.Findings {
		finding := map[string]interface{}{
			"category":    f.Category,
			"severity":    f.Severity,
			"description": f.Description,
			"suggestion":  f.Suggestion,
		}

		if f.Line != nil {
			finding["line"] = *f.Line
		}

		if f.FilePath != "" {
			finding["file_path"] = f.FilePath
		}

		if f.CodeSnippet != "" {
			finding["code_snippet"] = f.CodeSnippet
		}

		findings[i] = finding
	}

	result := map[string]interface{}{
		"findings":    findings,
		"summary":     resp.Summary,
		"provider":    resp.Provider,
		"duration_ms": resp.Duration.Milliseconds(),
	}

	if resp.Metadata != nil {
		metadata := map[string]interface{}{
			"source_type": resp.Metadata.SourceType,
		}

		if resp.Metadata.Model != "" {
			metadata["model"] = resp.Metadata.Model
		}
		if resp.Metadata.FileCount > 0 {
			metadata["file_count"] = resp.Metadata.FileCount
		}
		if resp.Metadata.LineCount > 0 {
			metadata["line_count"] = resp.Metadata.LineCount
		}
		if resp.Metadata.LinesAdded > 0 {
			metadata["lines_added"] = resp.Metadata.LinesAdded
		}
		if resp.Metadata.LinesRemoved > 0 {
			metadata["lines_removed"] = resp.Metadata.LinesRemoved
		}

		result["metadata"] = metadata
	}

	return result
}

// handleGitReview is a helper function for git-based review operations (staged, unstaged)
func (s *Server) handleGitReview(ctx context.Context, req *mcp.CallToolRequest, sourceType string) (*mcp.CallToolResult, error) {
	logging.Info(ctx, fmt.Sprintf("Handling review_%s request", sourceType))

	// Parse arguments
	var args struct {
		RepositoryPath string `json:"repository_path"`
		Provider       string `json:"provider,omitempty"`
		ReviewDepth    string `json:"review_depth,omitempty"`
	}

	if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	// Validate required arguments
	if args.RepositoryPath == "" {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: "repository_path is required"}},
		}, nil
	}

	if args.ReviewDepth == "" {
		args.ReviewDepth = "quick"
	}

	// Build review request
	reviewReq := review.Request{
		SourceType:     sourceType,
		RepositoryPath: args.RepositoryPath,
		Provider:       args.Provider,
		ReviewDepth:    args.ReviewDepth,
		Language:       "diff",
	}

	// Perform review (engine will populate Code from git)
	resp, err := s.engine.Review(ctx, reviewReq)
	if err != nil {
		logging.Error(ctx, "Review failed", "error", err)
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Review failed: %v", err)}},
		}, nil
	}

	// Format response as JSON content
	jsonData, err := json.MarshalIndent(formatReviewResponse(resp), "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to format response: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(jsonData)}},
	}, nil
}

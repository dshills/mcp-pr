package mcp

import (
	"context"
	"encoding/json"

	"github.com/dshills/mcp-pr/internal/review"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server represents the MCP code review server
type Server struct {
	mcpServer *mcp.Server
	engine    *review.Engine
}

// NewServer creates a new MCP server
func NewServer(engine *review.Engine) (*Server, error) {
	// Create MCP server implementation
	impl := &mcp.Implementation{
		Name:    "mcp-code-review",
		Version: "1.0.0",
	}

	opts := &mcp.ServerOptions{
		HasTools: true,
	}

	mcpServer := mcp.NewServer(impl, opts)

	srv := &Server{
		mcpServer: mcpServer,
		engine:    engine,
	}

	// Register tools
	srv.registerTools()

	return srv, nil
}

// registerTools registers all MCP tools
func (s *Server) registerTools() {
	// Register review_code tool
	s.mcpServer.AddTool(&mcp.Tool{
		Name:        "review_code",
		Description: "Review arbitrary code snippet for quality, security, and best practices",
		InputSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"code": {"type": "string", "description": "Code to review"},
				"language": {"type": "string", "description": "Programming language (e.g., go, python, javascript)"},
				"provider": {"type": "string", "enum": ["anthropic", "openai", "google"], "description": "LLM provider to use"},
				"review_depth": {"type": "string", "enum": ["quick", "thorough"], "default": "quick", "description": "Review depth"},
				"focus_areas": {"type": "array", "items": {"type": "string", "enum": ["bug", "security", "performance", "style", "best-practice"]}, "description": "Specific areas to focus on"}
			},
			"required": ["code", "language"]
		}`),
	}, s.handleReviewCode)

	// Register review_staged tool
	s.mcpServer.AddTool(&mcp.Tool{
		Name:        "review_staged",
		Description: "Review git staged changes in a repository",
		InputSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"repository_path": {"type": "string", "description": "Path to git repository"},
				"provider": {"type": "string", "enum": ["anthropic", "openai", "google"], "description": "LLM provider to use"},
				"review_depth": {"type": "string", "enum": ["quick", "thorough"], "default": "quick", "description": "Review depth"}
			},
			"required": ["repository_path"]
		}`),
	}, s.handleReviewStaged)

	// Register review_unstaged tool
	s.mcpServer.AddTool(&mcp.Tool{
		Name:        "review_unstaged",
		Description: "Review git unstaged changes in a repository",
		InputSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"repository_path": {"type": "string", "description": "Path to git repository"},
				"provider": {"type": "string", "enum": ["anthropic", "openai", "google"], "description": "LLM provider to use"},
				"review_depth": {"type": "string", "enum": ["quick", "thorough"], "default": "quick", "description": "Review depth"}
			},
			"required": ["repository_path"]
		}`),
	}, s.handleReviewUnstaged)

	// Register review_commit tool
	s.mcpServer.AddTool(&mcp.Tool{
		Name:        "review_commit",
		Description: "Review a specific git commit",
		InputSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"repository_path": {"type": "string", "description": "Path to git repository"},
				"commit_sha": {"type": "string", "description": "Git commit SHA to review"},
				"provider": {"type": "string", "enum": ["anthropic", "openai", "google"], "description": "LLM provider to use"},
				"review_depth": {"type": "string", "enum": ["quick", "thorough"], "default": "quick", "description": "Review depth"}
			},
			"required": ["repository_path", "commit_sha"]
		}`),
	}, s.handleReviewCommit)
}

// Run runs the server on stdin/stdout
func (s *Server) Run(ctx context.Context) error {
	return s.mcpServer.Run(ctx, &mcp.StdioTransport{})
}

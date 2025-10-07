# Research: MCP Code Review Server

**Feature**: MCP Code Review Server
**Phase**: 0 (Outline & Research)
**Date**: 2025-10-07

## Overview

This document consolidates research findings and technical decisions for implementing an MCP server that provides AI-powered code reviews using multiple LLM providers.

## Research Areas

### 1. MCP Server Implementation with go-sdk

**Decision**: Use github.com/modelcontextprotocol/go-sdk v1.0.0 for MCP server implementation

**Rationale**:
- Official Go SDK for Model Context Protocol maintained by the MCP team
- Provides stdio and SSE transport implementations
- Handles JSON-RPC 2.0 message framing and protocol compliance
- Includes tool and resource abstractions matching MCP specification
- Active development with v1.0.0 stable release

**Implementation Approach**:
- Create server using `mcp.NewServer()` with stdio transport (most common for MCP servers)
- Register tools for each review type: `review_code`, `review_staged`, `review_unstaged`, `review_commit`
- Use tool input schemas to validate parameters (code text, repository path, commit SHA, provider selection)
- Return structured JSON responses with review findings
- Handle errors through MCP error responses with appropriate error codes

**Alternatives Considered**:
- Build custom JSON-RPC 2.0 server: Rejected due to complexity and reinventing protocol handling
- Use generic RPC libraries: Rejected as they don't provide MCP-specific abstractions (tools, resources, prompts)

### 2. Multi-Provider LLM Integration

**Decision**: Implement provider adapter pattern with three concrete adapters

**Rationale**:
- Constitution Principle II requires provider-agnostic abstractions
- Each provider SDK has different APIs and authentication mechanisms
- Adapter pattern allows clean separation and testability
- Easy to add new providers in the future without changing core logic

**Provider Interface Design**:
```go
type Provider interface {
    Review(ctx context.Context, req ReviewRequest) (*ReviewResponse, error)
    Name() string
    IsAvailable() bool
}

type ReviewRequest struct {
    Code          string
    Language      string // e.g., "go", "python", "javascript"
    ReviewDepth   string // "quick" or "thorough"
    FocusAreas    []string // ["security", "performance", "style"]
}

type ReviewResponse struct {
    Findings []Finding
    Summary  string
    Duration time.Duration
}

type Finding struct {
    Category    string // "bug", "security", "performance", "style", "best-practice"
    Severity    string // "critical", "high", "medium", "low", "info"
    Line        *int   // Optional line number
    Description string
    Suggestion  string
}
```

**Provider-Specific Details**:

**Anthropic (Claude)**:
- SDK: github.com/anthropics/anthropic-sdk-go v1.13.0
- Auth: API key from ANTHROPIC_API_KEY environment variable
- Model: claude-3-5-sonnet-20241022 (latest recommended model)
- Prompt structure: System prompt with review instructions + user message with code
- Rate limits: 5 requests/minute (tier 1), handle with exponential backoff

**OpenAI (GPT)**:
- SDK: github.com/openai/openai-go/v3 v3.2.0
- Auth: API key from OPENAI_API_KEY environment variable
- Model: gpt-4-turbo or gpt-4o (configurable, default gpt-4o for speed)
- Prompt structure: System message + user message with code
- Rate limits: 500 requests/minute (paid tier), handle with retry logic

**Google (Gemini)**:
- SDK: google.golang.org/genai v1.28.0
- Auth: API key from GOOGLE_API_KEY environment variable
- Model: gemini-1.5-pro (best for code analysis)
- Prompt structure: Single text prompt with instructions and code
- Rate limits: 60 requests/minute (free tier), exponential backoff

**Prompt Engineering**:
All providers will use a structured prompt template:
```
You are a code review assistant. Analyze the following code and identify:
1. Bugs and logic errors
2. Security vulnerabilities
3. Performance issues
4. Style violations and best practice deviations

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
```language
<code_here>
```

Focus areas: <focus_areas>
Review depth: <quick|thorough>
```

**Alternatives Considered**:
- Single hardcoded provider: Violates Constitution Principle II
- Plugin architecture with dynamic loading: Over-engineered for current scope
- Provider chaining/fallback: Deferred to future enhancement

### 3. Git Integration

**Decision**: Use Go's os/exec package to shell out to git commands

**Rationale**:
- Git CLI is ubiquitous and stable
- Pure Go git libraries (go-git) have incomplete feature support for complex operations
- Direct git commands are well-documented and reliable
- Easy to test with mock exec functions

**Git Operations Needed**:
- `git diff --staged`: Get staged changes
- `git diff`: Get unstaged changes
- `git show <commit>`: Get specific commit diff
- `git rev-parse --verify <sha>`: Validate commit exists
- `git rev-parse --show-toplevel`: Get repository root
- `git diff --binary --no-color --unified=3`: Get diff with context

**Error Handling**:
- Check git is in PATH on startup
- Validate repository exists before operations
- Return structured errors for missing git, invalid repos, invalid commits
- Handle empty diffs gracefully (no changes to review)

**Diff Parsing**:
- Parse unified diff format to extract file paths and changed line numbers
- Skip binary files automatically
- Chunk large diffs (>5,000 lines) into smaller reviews if needed
- Preserve file context for multi-file diffs

**Alternatives Considered**:
- go-git library: Rejected due to incomplete feature set and complexity
- libgit2 with Go bindings: Rejected due to CGo dependency and platform issues

### 4. Configuration Management

**Decision**: Environment variables for all configuration

**Rationale**:
- MCP servers typically run as subprocesses; env vars are standard for configuration
- No need for config files (adds complexity and file I/O)
- API keys should never be in code or config files
- Simple to test with os.Setenv in tests

**Configuration Parameters**:
```
ANTHROPIC_API_KEY=<key>
OPENAI_API_KEY=<key>
GOOGLE_API_KEY=<key>

MCP_LOG_LEVEL=info|debug|warn|error (default: info)
MCP_DEFAULT_PROVIDER=anthropic|openai|google (default: anthropic)
MCP_REVIEW_TIMEOUT=30s (default: 30s)
MCP_MAX_DIFF_SIZE=10000 (lines, default: 10000)

# Per-provider timeouts
ANTHROPIC_TIMEOUT=30s
OPENAI_TIMEOUT=30s
GOOGLE_TIMEOUT=30s
```

**Alternatives Considered**:
- YAML/JSON config files: Over-engineered for simple configuration
- Command-line flags: Not standard for MCP servers (run as subprocesses)

### 5. Structured Logging

**Decision**: JSON logging using Go's slog package (standard library in Go 1.21+)

**Rationale**:
- Constitution Principle V requires structured logging
- slog is now standard library (no external dependencies)
- JSON format enables easy parsing and log aggregation
- Supports context-aware logging with fields
- Efficient and idiomatic Go

**Log Levels**:
- ERROR: API failures, git command failures, unrecoverable errors
- WARN: Timeouts, fallback provider usage, large diffs, missing git
- INFO: Review requests, provider selection, review completion, performance metrics
- DEBUG: Full request/response payloads, git command output, provider prompts

**Log Fields**:
- timestamp: ISO8601
- level: error|warn|info|debug
- message: Human-readable message
- operation: review_code|review_staged|review_unstaged|review_commit
- provider: anthropic|openai|google
- duration_ms: Operation duration
- error: Error message (if applicable)
- request_id: UUID for request tracing

**Alternatives Considered**:
- logrus: External dependency, unnecessary for simple structured logging
- zap: High-performance but overkill for MCP server use case
- Plain text logging: Violates Constitution Principle V

### 6. Testing Strategy

**Decision**: Three-tier testing per Constitution Principle IV

**Contract Tests**:
- MCP protocol compliance: Validate JSON-RPC 2.0 messages
- Provider interface: Ensure all adapters implement Provider interface correctly
- Tool schemas: Validate MCP tool input/output schemas

**Integration Tests**:
- Live API tests with each provider (requires API keys in CI)
- Git operations in temporary test repositories
- End-to-end MCP tool invocations with real providers

**Unit Tests**:
- Review engine logic with mocked providers
- Git client with mocked os/exec
- MCP tool handlers with fake requests
- Configuration loading

**Test Coverage Target**: ≥80% per Constitution

**Alternatives Considered**:
- Manual testing only: Violates Constitution Principle IV
- Integration tests only: Insufficient coverage of edge cases

## Unresolved Questions

None. All technical decisions finalized.

## Next Steps

Proceed to Phase 1: Generate data-model.md, contracts/, and quickstart.md based on these research findings.

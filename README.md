# MCP Code Review Server

[![CI](https://github.com/dshills/mcp-pr/actions/workflows/ci.yml/badge.svg)](https://github.com/dshills/mcp-pr/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dshills/mcp-pr)](https://goreportcard.com/report/github.com/dshills/mcp-pr)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/dshills/mcp-pr)](https://go.dev/doc/install)

An intelligent code review server implementing the [Model Context Protocol (MCP)](https://modelcontextprotocol.io) that provides AI-powered code analysis using Anthropic Claude, OpenAI GPT, or Google Gemini.

## Quick Start

```bash
# Install
go build -o mcp-code-review ./cmd/mcp-code-review

# Set API key (choose one or more)
export ANTHROPIC_API_KEY="your-anthropic-key"
export OPENAI_API_KEY="your-openai-key"
export GOOGLE_API_KEY="your-google-key"

# Run (MCP server mode)
./mcp-code-review
```

The server runs as an MCP subprocess, waiting for JSON-RPC requests over stdio.

---

## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
  - [Simple: Review Code Snippet](#simple-review-code-snippet)
  - [Git: Review Staged Changes](#git-review-staged-changes)
  - [Git: Review Unstaged Changes](#git-review-unstaged-changes)
  - [Git: Review Specific Commit](#git-review-specific-commit)
- [Tool Reference](#tool-reference)
- [Response Format](#response-format)
- [Examples](#examples)
- [Development](#development)
- [Architecture](#architecture)
- [Troubleshooting](#troubleshooting)

---

## Installation

### Prerequisites

- **Go 1.24+**: [Download Go](https://go.dev/dl/)
- **Git**: Required for git-based reviews
- **API Keys**: At least one LLM provider API key
  - [Anthropic Claude](https://console.anthropic.com/)
  - [OpenAI GPT](https://platform.openai.com/api-keys)
  - [Google Gemini](https://aistudio.google.com/app/apikey)

### Build from Source

```bash
# Clone repository
git clone https://github.com/dshills/mcp-pr.git
cd mcp-pr

# Build binary
go build -o mcp-code-review ./cmd/mcp-code-review

# Verify build
./mcp-code-review --version  # (if version flag supported)
```

### Install to GOPATH

```bash
go install github.com/dshills/mcp-pr/cmd/mcp-code-review@latest
```

The binary will be available at `$GOPATH/bin/mcp-code-review`.

---

## Configuration

### Required: API Keys

Set at least one provider API key as an environment variable:

```bash
# Anthropic Claude (recommended)
export ANTHROPIC_API_KEY="sk-ant-..."

# OpenAI GPT
export OPENAI_API_KEY="sk-..."

# Google Gemini
export GOOGLE_API_KEY="..."
```

### Optional: Environment Variables

```bash
# Logging
export MCP_PR_LOG_LEVEL=info          # debug|info|warn|error (default: info)

# Provider selection
export MCP_PR_DEFAULT_PROVIDER=anthropic  # anthropic|openai|google (default: anthropic)

# Timeouts (increased defaults for reliability)
export MCP_PR_REVIEW_TIMEOUT=120s     # Overall review timeout (default: 120s)
export ANTHROPIC_TIMEOUT=90s          # Anthropic API timeout (default: 90s)
export OPENAI_TIMEOUT=90s             # OpenAI API timeout (default: 90s)
export GOOGLE_TIMEOUT=90s             # Google API timeout (default: 90s)

# Diff size limits
export MCP_PR_MAX_DIFF_SIZE=10000     # Max diff size in bytes (default: 10000)
```

### MCP Client Configuration

If using Claude Desktop or another MCP client, add this server to your configuration:

**Claude Desktop** (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "code-review": {
      "command": "/path/to/mcp-code-review",
      "env": {
        "ANTHROPIC_API_KEY": "your-key-here"
      }
    }
  }
}
```

---

## Migration from v0.x

**Environment Variable Updates**: In recent versions, we renamed environment variables to be project-specific to avoid conflicts with other MCP servers.

### Variable Name Changes

| Old Name (Deprecated) | New Name | Status |
|-----------------------|----------|--------|
| `MCP_LOG_LEVEL` | `MCP_PR_LOG_LEVEL` | ⚠️ Deprecated, will be removed in v1.0.0 |
| `MCP_DEFAULT_PROVIDER` | `MCP_PR_DEFAULT_PROVIDER` | ⚠️ Deprecated, will be removed in v1.0.0 |
| `MCP_REVIEW_TIMEOUT` | `MCP_PR_REVIEW_TIMEOUT` | ⚠️ Deprecated, will be removed in v1.0.0 |
| `MCP_MAX_DIFF_SIZE` | `MCP_PR_MAX_DIFF_SIZE` | ⚠️ Deprecated, will be removed in v1.0.0 |

### What You Need to Do

**If you're upgrading from v0.x:**
1. Update your environment variable names from `MCP_*` to `MCP_PR_*`
2. Old names still work but will log deprecation warnings
3. New names take precedence if both are set

**What Stays the Same:**
- ✅ API key variables remain unchanged: `ANTHROPIC_API_KEY`, `OPENAI_API_KEY`, `GOOGLE_API_KEY`
- ✅ Provider-specific timeouts remain unchanged: `ANTHROPIC_TIMEOUT`, `OPENAI_TIMEOUT`, `GOOGLE_TIMEOUT`

**Example Migration:**

```bash
# Old configuration (still works with warnings)
export MCP_LOG_LEVEL=info
export MCP_DEFAULT_PROVIDER=anthropic

# New configuration (recommended)
export MCP_PR_LOG_LEVEL=info
export MCP_PR_DEFAULT_PROVIDER=anthropic
```

---

## Usage

The MCP server provides 4 tools for different code review workflows. Below are examples from simplest to most advanced.

### Simple: Review Code Snippet

Review any code snippet without git dependency.

**Tool**: `review_code`

**Example MCP Request**:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "review_code",
    "arguments": {
      "code": "func divide(a, b int) int {\n    return a / b\n}",
      "language": "go"
    }
  }
}
```

**Minimal Arguments**:
- `code` (required): Code to review
- `language` (optional): Programming language hint (e.g., "go", "python", "javascript")

**Optional Arguments**:
- `provider`: Choose LLM (`"anthropic"`, `"openai"`, `"google"`)
- `review_depth`: `"quick"` or `"thorough"` (default: `"quick"`)
- `focus_areas`: Array like `["security", "performance"]`

**Example Response**:

```json
{
  "findings": [
    {
      "category": "bug",
      "severity": "critical",
      "line": 2,
      "description": "Division by zero possible when b=0",
      "suggestion": "Add check: if b == 0 { return 0 } or return error"
    }
  ],
  "summary": "Found 1 critical issue: potential division by zero",
  "provider": "anthropic",
  "duration_ms": 1234,
  "metadata": {
    "source_type": "arbitrary",
    "model": "claude-3-7-sonnet-20250219"
  }
}
```

---

### Git: Review Staged Changes

Review code staged for commit (pre-commit workflow).

**Tool**: `review_staged`

**Use Case**: Run before `git commit` to catch issues early.

**Example Request**:

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "review_staged",
    "arguments": {
      "repository_path": "/path/to/your/repo"
    }
  }
}
```

**Arguments**:
- `repository_path` (required): Absolute path to git repository
- `provider` (optional): LLM provider to use
- `review_depth` (optional): `"quick"` or `"thorough"`

**What It Reviews**: Output of `git diff --staged` (all staged changes)

**Example Workflow**:

```bash
# Make changes
echo "new code" >> main.go

# Stage changes
git add main.go

# Review via MCP client (e.g., Claude Desktop)
# Tool: review_staged
# Args: { "repository_path": "/Users/me/myproject" }

# If review passes, commit
git commit -m "Add feature"
```

---

### Git: Review Unstaged Changes

Review working directory changes before staging (WIP feedback).

**Tool**: `review_unstaged`

**Use Case**: Get feedback on work-in-progress before staging.

**Example Request**:

```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "review_unstaged",
    "arguments": {
      "repository_path": "/path/to/your/repo",
      "review_depth": "quick"
    }
  }
}
```

**Arguments**: Same as `review_staged`

**What It Reviews**: Output of `git diff` (unstaged changes in working directory)

**Example Workflow**:

```bash
# Make changes (don't stage yet)
vim server.go

# Get quick feedback via MCP
# Tool: review_unstaged
# Args: { "repository_path": "/Users/me/myproject" }

# Fix issues, then stage
git add server.go
```

---

### Git: Review Specific Commit

Review a historical commit for audit or learning.

**Tool**: `review_commit`

**Use Case**: Understand what changed in a past commit, audit for security issues.

**Example Request**:

```json
{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "tools/call",
  "params": {
    "name": "review_commit",
    "arguments": {
      "repository_path": "/path/to/your/repo",
      "commit_sha": "a1b2c3d4e5f6",
      "provider": "openai",
      "review_depth": "thorough"
    }
  }
}
```

**Arguments**:
- `repository_path` (required): Absolute path to git repository
- `commit_sha` (required): Git commit SHA (full or short)
- `provider` (optional): LLM provider
- `review_depth` (optional): Review thoroughness

**What It Reviews**: Output of `git show <commit_sha>`

**Example Workflow**:

```bash
# Find suspicious commit
git log --oneline

# Review it via MCP
# Tool: review_commit
# Args: {
#   "repository_path": "/Users/me/myproject",
#   "commit_sha": "abc123"
# }

# Check for security issues in that commit
```

---

## Tool Reference

### `review_code`

Review arbitrary code text without git.

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `code` | string | ✅ | - | Code to review |
| `language` | string | ❌ | auto-detect | Language hint (go, python, js, etc.) |
| `provider` | string | ❌ | env default | `anthropic`, `openai`, or `google` |
| `review_depth` | string | ❌ | `quick` | `quick` or `thorough` |
| `focus_areas` | array | ❌ | all | `["bug", "security", "performance", "style", "best-practice"]` |

### `review_staged`

Review git staged changes (`git diff --staged`).

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `repository_path` | string | ✅ | - | Absolute path to git repository |
| `provider` | string | ❌ | env default | `anthropic`, `openai`, or `google` |
| `review_depth` | string | ❌ | `quick` | `quick` or `thorough` |

### `review_unstaged`

Review git unstaged changes (`git diff`).

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `repository_path` | string | ✅ | - | Absolute path to git repository |
| `provider` | string | ❌ | env default | `anthropic`, `openai`, or `google` |
| `review_depth` | string | ❌ | `quick` | `quick` or `thorough` |

### `review_commit`

Review specific commit (`git show <sha>`).

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `repository_path` | string | ✅ | - | Absolute path to git repository |
| `commit_sha` | string | ✅ | - | Git commit SHA (full or short) |
| `provider` | string | ❌ | env default | `anthropic`, `openai`, or `google` |
| `review_depth` | string | ❌ | `quick` | `quick` or `thorough` |

---

## Response Format

All tools return JSON with this structure:

```typescript
{
  findings: Array<{
    category: "bug" | "security" | "performance" | "style" | "best-practice",
    severity: "critical" | "high" | "medium" | "low" | "info",
    line?: number,              // Line number (if applicable)
    file_path?: string,         // File path (for git reviews)
    description: string,        // What the issue is
    suggestion: string,         // How to fix it
    code_snippet?: string       // Relevant code excerpt
  }>,
  summary: string,              // Overall assessment
  provider: string,             // Which LLM was used
  duration_ms: number,          // Review duration
  metadata: {
    source_type: string,        // "arbitrary", "staged", "unstaged", "commit"
    model: string,              // LLM model name
    file_count?: number,        // Number of files (git reviews)
    line_count?: number,        // Total lines reviewed
    lines_added?: number,       // Lines added (git diffs)
    lines_removed?: number      // Lines removed (git diffs)
  }
}
```

### Finding Categories

- **bug**: Logic errors, crashes, incorrect behavior
- **security**: Vulnerabilities, injection risks, auth issues
- **performance**: Inefficiencies, memory leaks, slow algorithms
- **style**: Code formatting, naming conventions, readability
- **best-practice**: Design patterns, maintainability, testability

### Severity Levels

- **critical**: Must fix immediately (security hole, crash)
- **high**: Should fix soon (data loss risk, major bug)
- **medium**: Fix when convenient (performance issue)
- **low**: Nice to have (minor style issue)
- **info**: Informational (suggestion, tip)

---

## Examples

### Example 1: Review Python Function

**Code**:
```python
def fetch_user(user_id):
    query = f"SELECT * FROM users WHERE id = {user_id}"
    return db.execute(query)
```

**MCP Request**:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "review_code",
    "arguments": {
      "code": "def fetch_user(user_id):\n    query = f\"SELECT * FROM users WHERE id = {user_id}\"\n    return db.execute(query)",
      "language": "python",
      "focus_areas": ["security"]
    }
  }
}
```

**Expected Findings**:
- **Security/Critical**: SQL injection vulnerability
- **Suggestion**: Use parameterized queries

---

### Example 2: Pre-Commit Hook Integration

Create a git hook to review staged changes before committing:

**`.git/hooks/pre-commit`**:
```bash
#!/bin/bash

# Review staged changes via MCP
REPO_PATH=$(git rev-parse --show-toplevel)

# Use MCP client to call review_staged
# (Implementation depends on your MCP client)

# Example with Claude Desktop CLI (if available)
mcp-client call code-review review_staged \
  --repository_path "$REPO_PATH" \
  --review_depth quick

# Exit with error if critical issues found
if [ $? -ne 0 ]; then
  echo "❌ Code review found critical issues. Fix before committing."
  exit 1
fi
```

Make it executable:
```bash
chmod +x .git/hooks/pre-commit
```

---

### Example 3: CI/CD Integration

Review commits in GitHub Actions:

**`.github/workflows/code-review.yml`**:
```yaml
name: AI Code Review

on: [pull_request]

jobs:
  review:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.24'

      - name: Install MCP Code Review
        run: |
          go install github.com/dshills/mcp-pr/cmd/mcp-code-review@latest

      - name: Review Pull Request Changes
        env:
          ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
        run: |
          # Get changed files and review them
          git diff origin/main...HEAD > changes.diff

          # Call MCP server to review
          # (Implement MCP client call here)
```

---

## Development

### Build

```bash
# Development build
go build -o mcp-code-review ./cmd/mcp-code-review

# Production build with optimizations
go build -ldflags="-s -w" -o mcp-code-review ./cmd/mcp-code-review
```

### Test

```bash
# Run all tests
go test ./tests/...

# Run with coverage
go test -cover ./...

# Run specific test suite
go test ./tests/contract/...     # Contract tests
go test ./tests/integration/...  # Integration tests (requires API keys)
go test ./tests/unit/...         # Unit tests

# Verbose output
go test -v ./tests/...
```

### Integration Tests

Integration tests require API keys:

```bash
export ANTHROPIC_API_KEY="your-key"
export OPENAI_API_KEY="your-key"
export GOOGLE_API_KEY="your-key"

go test ./tests/integration/... -v
```

### Linting

```bash
# Run linter
golangci-lint run ./...

# Fix auto-fixable issues
golangci-lint run --fix ./...
```

### Code Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View HTML report
go tool cover -html=coverage.out

# Target: ≥80% (per project constitution)
```

---

## Architecture

### Directory Structure

```
mcp-pr/
├── cmd/
│   └── mcp-code-review/        # Main server entry point
│       └── main.go
├── internal/
│   ├── config/                 # Configuration loading
│   │   └── config.go
│   ├── git/                    # Git operations
│   │   ├── client.go           # Git command wrappers
│   │   └── diff.go             # Diff parsing
│   ├── logging/                # Structured logging
│   │   └── logger.go
│   ├── mcp/                    # MCP protocol
│   │   ├── server.go           # Server initialization
│   │   └── tools.go            # Tool handlers
│   ├── providers/              # LLM providers
│   │   ├── provider.go         # Provider interface
│   │   ├── anthropic.go        # Claude integration
│   │   ├── openai.go           # GPT integration
│   │   └── google.go           # Gemini integration
│   └── review/                 # Review engine
│       ├── engine.go           # Review orchestration
│       ├── request.go          # Request models
│       └── response.go         # Response models
├── tests/
│   ├── contract/               # Contract tests
│   ├── integration/            # Integration tests
│   └── unit/                   # Unit tests
└── go.mod
```

### Components

#### MCP Server (`internal/mcp/`)
- Implements JSON-RPC 2.0 over stdio
- Registers 4 tools (review_code, review_staged, review_unstaged, review_commit)
- Handles request parsing and response formatting

#### Review Engine (`internal/review/`)
- Orchestrates review workflow
- Manages provider selection and failover
- Implements retry logic (1 retry = 2 total attempts)
- Validates diff size limits before sending to LLM
- Populates code from git for git-based reviews
- Progress logging throughout review lifecycle

#### Git Client (`internal/git/`)
- Wraps git commands (diff, show, rev-parse) with context-aware timeouts
- 30s timeout for diff operations (prevents hanging on slow filesystems)
- 10s timeout for validation operations
- Parses unified diff format
- Extracts file paths, line numbers, and changes

#### Providers (`internal/providers/`)
- Common interface: `Review(ctx, Request) (*Response, error)`
- Anthropic: Claude 3.7 Sonnet
- OpenAI: GPT-4o
- Google: Gemini (deprecated SDK, may require migration)

---

## Troubleshooting

### "No API key found"

**Problem**: Server exits with "at least one provider API key required"

**Solution**: Set at least one API key:
```bash
export ANTHROPIC_API_KEY="sk-ant-..."
```

---

### "Not a git repository"

**Problem**: Git tools fail with "fatal: not a git directory"

**Solution**: Ensure `repository_path` points to a valid git repository:
```bash
# Check if directory is a git repo
cd /path/to/repo
git rev-parse --git-dir  # Should output .git

# If not, initialize it
git init
```

---

### "Provider timeout" or "Operation timed out"

**Problem**: Review fails with timeout error after 90-120 seconds

**Solution**: The server now has improved timeout handling with these defaults:
- Provider API timeouts: **90s** (up from 30s)
- Git operation timeouts: **30s** (prevents hanging on slow filesystems)
- Overall review timeout: **120s** (up from 30s)

If you still experience timeouts, you can increase them:
```bash
# Increase provider timeout to 180s
export ANTHROPIC_TIMEOUT=180s
export OPENAI_TIMEOUT=180s
export GOOGLE_TIMEOUT=180s

# Or increase overall review timeout
export MCP_REVIEW_TIMEOUT=300s
```

**Alternative solutions**:
1. Use `quick` review depth for faster responses (default)
2. Review smaller diffs (break large changes into smaller commits)
3. Increase `MCP_MAX_DIFF_SIZE` if hitting size limits:
   ```bash
   export MCP_MAX_DIFF_SIZE=50000  # Increase from default 10000 bytes
   ```

**Note**: Git operations (diff, show) now have built-in 30s timeouts to prevent indefinite hangs on network filesystems or large repositories.

---

### "Diff size exceeds maximum allowed size"

**Problem**: Review fails with "diff size (X bytes) exceeds maximum allowed size (10000 bytes)"

**Solution**: This is a safety limit to prevent timeouts on very large diffs. You have three options:

1. **Increase the limit** (recommended for legitimate large diffs):
   ```bash
   export MCP_MAX_DIFF_SIZE=50000  # Increase to 50KB
   ```

2. **Break changes into smaller commits**:
   ```bash
   git add -p  # Stage changes interactively in smaller chunks
   git commit -m "Part 1: ..."
   ```

3. **Review specific files** separately using `review_code` tool instead of git-based reviews

**Note**: Very large diffs (>50KB) may still cause timeouts even with increased limits. Consider reviewing critical files separately.

---

### "Invalid commit SHA"

**Problem**: `review_commit` fails with "invalid commit SHA"

**Solution**: Verify commit exists:
```bash
git rev-parse --verify abc123  # Replace with your SHA
```

Use full or short SHA (minimum 7 characters):
```json
{
  "commit_sha": "a1b2c3d"  // Short SHA (7+ chars)
}
```

---

### "Empty diff / No changes to review"

**Problem**: Git reviews return empty results

**Solution**: Verify there are changes:
```bash
# For staged review
git diff --staged  # Should show changes

# For unstaged review
git diff  # Should show changes

# If no output, there are no changes to review
```

---

### Integration Tests Fail

**Problem**: Tests fail with API errors

**Solution**: Ensure valid API keys are set:
```bash
export ANTHROPIC_API_KEY="valid-key"
export OPENAI_API_KEY="valid-key"

# Skip integration tests if keys unavailable
go test ./tests/unit/...      # Run only unit tests
go test ./tests/contract/...  # Run only contract tests
```

---

## Performance Tips

### Timeout Improvements (v1.1.0+)

The server now includes several optimizations to prevent timeouts:

- **Increased default timeouts**: Provider APIs now have 90s (up from 30s)
- **Git operation timeouts**: 30s limit prevents indefinite hangs on slow filesystems
- **Reduced retry attempts**: 1 retry (down from 3) minimizes worst-case delays
- **Diff size validation**: 10KB default limit prevents oversized requests
- **Progress logging**: Visibility into review stages for debugging

### Best Practices

1. **Use `quick` review depth** for faster feedback (default)
2. **Choose faster providers**: Anthropic Claude is generally fastest
3. **Review smaller diffs**: Break large changes into smaller commits
   ```bash
   # Stage files selectively
   git add file1.go file2.go
   git commit -m "Part 1: Core logic"

   git add file3.go file4.go
   git commit -m "Part 2: Tests"
   ```
4. **Monitor diff sizes**: Keep diffs under 10KB for best performance
5. **Increase timeouts for complex reviews**: Use `thorough` depth with longer timeouts:
   ```bash
   export ANTHROPIC_TIMEOUT=180s
   export MCP_REVIEW_TIMEOUT=300s
   ```

---

## Security Considerations

- **API Keys**: Never commit API keys to version control
- **Code Privacy**: Code is sent to third-party LLM providers (Anthropic/OpenAI/Google)
- **Local Git**: Git operations are local; diffs are not sent anywhere except to the LLM
- **Secrets in Code**: Review may detect hardcoded secrets, but don't rely on it exclusively

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines.

### Running Tests

```bash
# All tests
go test ./tests/... -v

# With coverage
go test -cover ./...
```

### Code Quality

```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Lint
golangci-lint run ./...
```

---

## License

See [LICENSE](LICENSE) file for details.

---

## Links

- [Model Context Protocol](https://modelcontextprotocol.io)
- [Anthropic Claude](https://www.anthropic.com/claude)
- [OpenAI GPT](https://openai.com/gpt-4)
- [Google Gemini](https://ai.google.dev/)

---

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Code of Conduct

This project adheres to the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

### Security

For security vulnerabilities, please see our [Security Policy](SECURITY.md).

---

## Support

For issues, questions, or feature requests, please [open an issue](https://github.com/dshills/mcp-pr/issues/new/choose) on GitHub.

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Built with ❤️ using Go and the Model Context Protocol**

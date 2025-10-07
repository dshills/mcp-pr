# Quickstart: MCP Code Review Server

**Feature**: MCP Code Review Server
**Phase**: 1 (Design & Contracts)
**Date**: 2025-10-07

## Purpose

This quickstart guide provides integration test scenarios and usage examples for validating the MCP Code Review Server implementation.

## Prerequisites

- Go 1.25.1+
- Git installed and in PATH
- At least one LLM provider API key configured:
  - `ANTHROPIC_API_KEY` for Anthropic Claude
  - `OPENAI_API_KEY` for OpenAI GPT
  - `GOOGLE_API_KEY` for Google Gemini

## Installation

```bash
# Build the MCP server
go build -o mcp-code-review ./cmd/mcp-code-review

# Or install to GOPATH
go install ./cmd/mcp-code-review
```

## Quick Test: Review Arbitrary Code

```bash
# Start the MCP server (stdio transport)
./mcp-code-review

# In another terminal, send MCP request (JSON-RPC 2.0)
echo '{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "review_code",
    "arguments": {
      "code": "func divide(a, b int) int { return a / b }",
      "language": "go",
      "provider": "anthropic",
      "review_depth": "quick"
    }
  }
}' | ./mcp-code-review
```

**Expected Response**:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "findings": [
      {
        "category": "bug",
        "severity": "high",
        "line": 1,
        "description": "Division by zero vulnerability: function does not handle b == 0",
        "suggestion": "Add validation: if b == 0 { return error or panic }"
      }
    ],
    "summary": "Found 1 high-severity bug: potential division by zero",
    "provider": "anthropic",
    "duration_ms": 1234
  }
}
```

## Integration Test Scenarios

### Scenario 1: Review Arbitrary Code (P1 User Story)

**Objective**: Verify basic code review functionality without git dependency

**Setup**:
```bash
export ANTHROPIC_API_KEY="your-key-here"
```

**Test Case 1.1**: Review code with security vulnerability
```bash
# Input: SQL injection vulnerability
CODE='db.Query("SELECT * FROM users WHERE name = " + userInput)'

# Expected: Security finding with "critical" or "high" severity
# Expected: Suggestion to use parameterized queries
```

**Test Case 1.2**: Review clean code
```bash
# Input: Well-written code
CODE='func add(a, b int) int { return a + b }'

# Expected: Empty findings array
# Expected: Summary: "No issues found" or similar positive message
```

**Test Case 1.3**: Empty code validation
```bash
# Input: Empty string
CODE=''

# Expected: Error response with validation failure message
```

**Test Case 1.4**: Multi-provider comparison
```bash
# Run same code review with all three providers
# Input: Same code snippet
# Providers: anthropic, openai, google

# Expected: All return findings (may differ in wording but similar severity)
# Expected: Response times < 10 seconds
```

---

### Scenario 2: Review Git Staged Changes (P2 User Story)

**Objective**: Verify git staged diff review workflow

**Setup**:
```bash
# Create test repository
mkdir /tmp/test-repo
cd /tmp/test-repo
git init
echo "func foo() {}" > main.go
git add main.go
git commit -m "Initial commit"

# Make changes and stage them
echo "func bar() { panic('test') }" >> main.go
git add main.go
```

**Test Case 2.1**: Review staged changes
```bash
# MCP request
{
  "method": "tools/call",
  "params": {
    "name": "review_staged",
    "arguments": {
      "repository_path": "/tmp/test-repo",
      "provider": "anthropic"
    }
  }
}

# Expected: Findings related to panic usage
# Expected: metadata.source_type = "staged"
# Expected: metadata.file_count = 1
# Expected: metadata.lines_added > 0
```

**Test Case 2.2**: No staged changes
```bash
# Setup: No staged files
git reset HEAD main.go

# Expected: Success response with empty findings
# Expected: Summary indicating nothing to review
```

**Test Case 2.3**: Multiple file staged review
```bash
# Stage multiple files with different issues
echo "var x = y / 0" > bug.go
echo "password := 'hardcoded'" > security.go
git add bug.go security.go

# Expected: Findings for both files
# Expected: metadata.file_count = 2
# Expected: Each finding has file_path set correctly
```

---

### Scenario 3: Review Unstaged Changes (P3 User Story)

**Objective**: Verify working directory diff review

**Setup**:
```bash
cd /tmp/test-repo
# Modify without staging
echo "func baz() int { return 1/0 }" >> main.go
```

**Test Case 3.1**: Review unstaged modifications
```bash
{
  "method": "tools/call",
  "params": {
    "name": "review_unstaged",
    "arguments": {
      "repository_path": "/tmp/test-repo"
    }
  }
}

# Expected: Finding about division by zero
# Expected: metadata.source_type = "unstaged"
```

**Test Case 3.2**: Untracked files handling
```bash
# Create new untracked file
echo "func new() {}" > new.go

# Expected: System either includes untracked files or clearly indicates they're skipped
# Expected: Documented behavior in response or logs
```

---

### Scenario 4: Review Specific Commit (P4 User Story)

**Objective**: Verify historical commit review

**Setup**:
```bash
cd /tmp/test-repo
git log --oneline  # Get commit SHA
```

**Test Case 4.1**: Review valid commit
```bash
{
  "method": "tools/call",
  "params": {
    "name": "review_commit",
    "arguments": {
      "repository_path": "/tmp/test-repo",
      "commit_sha": "abc1234",  # Use actual SHA
      "provider": "openai"
    }
  }
}

# Expected: Review of that commit's changes
# Expected: metadata.source_type = "commit"
```

**Test Case 4.2**: Invalid commit SHA
```bash
{
  "arguments": {
    "commit_sha": "invalid"
  }
}

# Expected: Error response
# Expected: Clear message indicating invalid SHA
```

**Test Case 4.3**: Merge commit
```bash
# Create merge commit
git checkout -b feature
echo "func feature() {}" > feature.go
git add feature.go && git commit -m "Add feature"
git checkout main
git merge feature

# Review merge commit
# Expected: Reviews all changes introduced by merge
```

---

## Edge Case Testing

### Edge Case 1: Large Diff Handling
```bash
# Generate large diff (>5000 lines)
for i in {1..6000}; do echo "// Line $i" >> large.go; done
git add large.go

# Review with review_staged
# Expected: Chunking or warning about size
# Expected: Response within timeout (30s)
# Expected: metadata.line_count accurately reflects size
```

### Edge Case 2: Binary Files
```bash
# Add binary file
echo "binary" | gzip > binary.gz
git add binary.gz

# Expected: Binary files skipped
# Expected: Clear indication in logs or metadata
```

### Edge Case 3: API Failure Handling
```bash
# Set invalid API key
export ANTHROPIC_API_KEY="invalid"

# Attempt review
# Expected: Error response with clear message
# Expected: Retry logic triggered (check logs)
# Expected: If retries exhausted, final error within 5s (per SC-007)
```

### Edge Case 4: Git Not in PATH
```bash
# Temporarily remove git from PATH
export PATH="/usr/bin"  # Minimal path without git

# Attempt git-based review
# Expected: Clear error message indicating git not found
# Expected: Suggestion to install git or add to PATH
```

### Edge Case 5: Repository Without Commits
```bash
# Fresh git init
mkdir /tmp/empty-repo
cd /tmp/empty-repo
git init

# Attempt any git-based review
# Expected: Graceful handling with appropriate error message
```

### Edge Case 6: Mixed Line Endings
```bash
# Create file with CRLF
printf "func test() {\r\n  return true\r\n}" > crlf.go
git add crlf.go

# Expected: Review works regardless of line endings
# Expected: Line numbers accurate
```

---

## Performance Validation

Based on Success Criteria from spec.md:

### SC-001: Review latency <10s for 500-line snippets
```bash
# Test with 500-line code snippet
# Measure: duration_ms in response
# Expected: <10,000 ms
```

### SC-002: 95% success rate across providers
```bash
# Run 100 reviews split across 3 providers
# Expected: ≥95 successful responses
```

### SC-005: Handle up to 100 changed files
```bash
# Create 100 files with small changes
# Stage and review
# Expected: Completes without timeout
```

### SC-007: Provider failure recovery <5s
```bash
# Simulate failure (invalid key)
# Measure time from request to error response
# Expected: <5000 ms total (including retries)
```

---

## Manual Testing Checklist

- [ ] MCP server starts without errors
- [ ] All 4 tools registered and discoverable via tools/list
- [ ] review_code works with all 3 providers
- [ ] review_staged identifies issues in staged changes
- [ ] review_unstaged identifies issues in working directory
- [ ] review_commit works with valid SHA
- [ ] Invalid inputs return clear error messages
- [ ] Structured logging outputs JSON format
- [ ] All log levels work (ERROR, WARN, INFO, DEBUG)
- [ ] Configuration via environment variables works
- [ ] Empty diffs handled gracefully
- [ ] Binary files skipped correctly
- [ ] Large diffs (>5000 lines) handled
- [ ] Provider API failures trigger retries
- [ ] Git command failures logged with context

---

## Troubleshooting

### Issue: "git command not found"
**Solution**: Ensure git is installed and in PATH

### Issue: "Provider API key not configured"
**Solution**: Set ANTHROPIC_API_KEY, OPENAI_API_KEY, or GOOGLE_API_KEY

### Issue: "Repository not found"
**Solution**: Verify repository_path is absolute and contains .git directory

### Issue: Timeout on large reviews
**Solution**: Increase MCP_REVIEW_TIMEOUT or chunk the diff

### Issue: Findings don't include line numbers
**Solution**: Check diff parsing logic for line number extraction

---

## Next Steps

After validation:
1. Run automated test suite: `go test ./tests/...`
2. Check code coverage: `go test -cover ./...` (target: ≥80%)
3. Run linters: `gofmt`, `go vet`, `golangci-lint`
4. Deploy MCP server to production environment
5. Monitor structured logs for errors and performance metrics

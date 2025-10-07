# MCP Tool Contracts

This directory contains the formal API contracts for MCP Code Review Server tools.

## Files

- **mcp-tools.json**: JSON Schema definitions for all MCP tools, including input/output schemas

## Tools

### 1. review_code
Review arbitrary code text without git integration.

**Input**: code, language (optional), provider, review_depth, focus_areas
**Output**: findings[], summary, provider, duration_ms, metadata

### 2. review_staged
Review git staged changes (pre-commit workflow).

**Input**: repository_path (optional), provider, review_depth, focus_areas
**Output**: findings[], summary, provider, duration_ms, metadata

### 3. review_unstaged
Review git unstaged changes (work-in-progress feedback).

**Input**: repository_path (optional), provider, review_depth, focus_areas
**Output**: findings[], summary, provider, duration_ms, metadata

### 4. review_commit
Review a specific commit by SHA (audit/learning).

**Input**: repository_path, commit_sha, provider, review_depth, focus_areas
**Output**: findings[], summary, provider, duration_ms, metadata

## Shared Types

### Finding
```json
{
  "category": "bug|security|performance|style|best-practice",
  "severity": "critical|high|medium|low|info",
  "line": 42 or null,
  "file_path": "path/to/file.go",
  "description": "Issue explanation",
  "suggestion": "How to fix it",
  "code_snippet": "relevant code"
}
```

### ReviewMetadata
```json
{
  "source_type": "arbitrary|staged|unstaged|commit",
  "file_count": 3,
  "line_count": 150,
  "lines_added": 75,
  "lines_removed": 25,
  "model": "claude-3-5-sonnet-20241022"
}
```

## Providers

- **anthropic**: Claude models (default: claude-3-5-sonnet-20241022)
- **openai**: GPT models (default: gpt-4o)
- **google**: Gemini models (default: gemini-1.5-pro)

## Review Depths

- **quick**: Fast scan for obvious issues (default)
- **thorough**: Deep analysis including subtle code smells

## Focus Areas

- **bugs**: Logic errors, null pointer issues, off-by-one errors
- **security**: SQL injection, XSS, hardcoded secrets, insecure crypto
- **performance**: O(n²) algorithms, memory leaks, inefficient queries
- **style**: Naming conventions, formatting, code organization
- **best-practices**: Language idioms, design patterns, maintainability

If focus_areas is empty, all areas are reviewed.

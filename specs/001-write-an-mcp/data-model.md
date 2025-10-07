# Data Model: MCP Code Review Server

**Feature**: MCP Code Review Server
**Phase**: 1 (Design & Contracts)
**Date**: 2025-10-07

## Overview

This document defines the core data entities and their relationships for the MCP Code Review Server. The system is stateless with no persistent storage; all entities are runtime request/response models.

## Core Entities

### 1. ReviewRequest

Represents an incoming request for code review.

**Fields**:
- `SourceType` (string, required): One of "arbitrary", "staged", "unstaged", "commit"
- `Code` (string, optional): Raw code text (required if SourceType="arbitrary", otherwise empty)
- `RepositoryPath` (string, optional): Path to git repository (required if SourceType != "arbitrary")
- `CommitSHA` (string, optional): Git commit SHA (required if SourceType="commit")
- `Provider` (string, required): One of "anthropic", "openai", "google"
- `Language` (string, optional): Programming language hint (e.g., "go", "python"); auto-detected if empty
- `ReviewDepth` (string, optional): "quick" or "thorough"; defaults to "quick"
- `FocusAreas` ([]string, optional): Filter review to specific categories; empty = all categories

**Validation Rules**:
- If SourceType="arbitrary", Code must be non-empty
- If SourceType="staged|unstaged|commit", RepositoryPath must be valid directory
- If SourceType="commit", CommitSHA must be non-empty
- Provider must be one of the supported values
- ReviewDepth must be "quick" or "thorough" if provided
- FocusAreas values must be in ["bugs", "security", "performance", "style", "best-practices"]

**Relationships**:
- Input to Provider.Review() method
- Constructed from MCP tool invocation parameters

---

### 2. ReviewResponse

Contains the results of a code review.

**Fields**:
- `Findings` ([]Finding, required): Array of identified issues (may be empty)
- `Summary` (string, required): Overall assessment paragraph
- `Provider` (string, required): Provider that generated this review
- `Duration` (duration, required): Time taken to generate review
- `Metadata` (ReviewMetadata, optional): Additional context about the review

**Validation Rules**:
- Findings array can be empty (perfect code)
- Summary must be non-empty (even for clean code: "No issues found")
- Duration must be positive
- Provider must match request provider

**Relationships**:
- Returned from Provider.Review() method
- Serialized as MCP tool response

---

### 3. Finding

Represents a single issue identified in code review.

**Fields**:
- `Category` (string, required): One of "bug", "security", "performance", "style", "best-practice"
- `Severity` (string, required): One of "critical", "high", "medium", "low", "info"
- `Line` (*int, optional): Line number where issue occurs (null for file-level issues)
- `FilePath` (string, optional): Relative file path (for git diffs with multiple files)
- `Description` (string, required): Clear explanation of the issue
- `Suggestion` (string, required): Actionable remediation advice
- `CodeSnippet` (string, optional): Relevant code excerpt

**Validation Rules**:
- Category must be one of the enum values
- Severity must be one of the enum values
- Line, if provided, must be positive
- Description and Suggestion must be non-empty
- CodeSnippet, if provided, should be <100 characters for brevity

**Relationships**:
- Component of ReviewResponse
- Parsed from provider LLM JSON output

**Severity Priority**:
- critical: Security vulnerabilities, data loss bugs
- high: Logic errors, performance bottlenecks
- medium: Code smells, minor bugs
- low: Style inconsistencies, minor optimizations
- info: Suggestions, educational notes

---

### 4. ReviewMetadata

Optional metadata about the review context.

**Fields**:
- `SourceType` (string, required): Echo of request source type
- `FileCount` (int, optional): Number of files reviewed (for git diffs)
- `LineCount` (int, optional): Total lines of code reviewed
- `LinesAdded` (int, optional): Lines added (for git diffs)
- `LinesRemoved` (int, optional): Lines removed (for git diffs)
- `Model` (string, optional): Specific LLM model used (e.g., "claude-3-5-sonnet-20241022")

**Relationships**:
- Nested in ReviewResponse
- Provides context for review consumers

---

### 5. ProviderConfig

Configuration for an LLM provider adapter.

**Fields**:
- `Name` (string, required): Provider name ("anthropic", "openai", "google")
- `APIKey` (string, required): API authentication key
- `Model` (string, required): Model identifier
- `Timeout` (duration, required): API request timeout
- `MaxRetries` (int, required): Number of retry attempts
- `RetryDelay` (duration, required): Delay between retries

**Validation Rules**:
- Name must be non-empty
- APIKey must be non-empty
- Model must be non-empty
- Timeout must be >0
- MaxRetries must be ≥0
- RetryDelay must be ≥0

**Relationships**:
- Loaded from environment variables at server startup
- Passed to provider adapters during initialization

---

### 6. GitDiff

Represents a git diff with metadata.

**Fields**:
- `Content` (string, required): Raw unified diff output
- `Files` ([]DiffFile, required): Parsed file-level changes
- `Stats` (DiffStats, required): Summary statistics

**Validation Rules**:
- Content must be non-empty (caller should check for empty diffs before creating)
- Files array must have ≥1 entry
- Stats must accurately reflect files

**Relationships**:
- Generated by git.Client
- Input to review engine

---

### 7. DiffFile

Represents changes in a single file within a diff.

**Fields**:
- `Path` (string, required): Relative file path
- `ChangeType` (string, required): "added", "modified", "deleted", "renamed"
- `LinesAdded` (int, required): Count of added lines
- `LinesRemoved` (int, required): Count of removed lines
- `Chunks` ([]DiffChunk, required): Individual diff hunks

**Validation Rules**:
- Path must be non-empty
- ChangeType must be enum value
- LinesAdded/LinesRemoved must be ≥0
- Chunks array can be empty for binary files

**Relationships**:
- Component of GitDiff
- Corresponds to one file's changes in a git diff

---

### 8. DiffChunk

Represents a single hunk in a unified diff.

**Fields**:
- `OldStart` (int, required): Starting line in old file
- `OldCount` (int, required): Number of lines in old file
- `NewStart` (int, required): Starting line in new file
- `NewCount` (int, required): Number of lines in new file
- `Lines` ([]DiffLine, required): Individual line changes

**Validation Rules**:
- All integer fields must be ≥0
- Lines array must have ≥1 entry

**Relationships**:
- Component of DiffFile
- Parsed from unified diff @@ headers

---

### 9. DiffLine

Represents a single line in a diff chunk.

**Fields**:
- `Type` (string, required): "context", "added", "removed"
- `Content` (string, required): Line content (without leading +/-)
- `LineNumber` (*int, optional): Line number in new file (null for removed lines)

**Validation Rules**:
- Type must be enum value
- Content can be empty (blank lines)
- LineNumber must be positive if provided

**Relationships**:
- Component of DiffChunk
- Used to associate findings with specific line numbers

---

### 10. MCPTool

Represents an MCP tool definition.

**Fields**:
- `Name` (string, required): Tool name (e.g., "review_code")
- `Description` (string, required): Human-readable tool description
- `InputSchema` (JSONSchema, required): JSON schema for tool parameters

**Examples**:
- `review_code`: Review arbitrary code text
- `review_staged`: Review git staged changes
- `review_unstaged`: Review git unstaged changes
- `review_commit`: Review a specific git commit

**Relationships**:
- Registered with MCP server
- Maps to ReviewRequest factory functions

---

## Entity Relationships Diagram

```
MCPTool (4 tools)
    ↓ invoked
ReviewRequest
    ↓ passed to
Provider.Review()
    ↓ calls
LLM API (Anthropic/OpenAI/Google)
    ↓ returns
ReviewResponse
    ├─ contains → []Finding
    └─ contains → ReviewMetadata

(Parallel flow for git-based reviews)
ReviewRequest (source_type != "arbitrary")
    ↓ triggers
GitClient.GetDiff()
    ↓ returns
GitDiff
    ├─ contains → []DiffFile
    │     ├─ contains → []DiffChunk
    │     │     └─ contains → []DiffLine
    │     └─ contains → DiffStats
    ↓ passed to
Provider.Review() [as code input]
```

## State Transitions

### Review Lifecycle

1. **Created**: ReviewRequest constructed from MCP tool invocation
2. **Fetching** (if git-based): Git client retrieves diff
3. **Validating**: Request validation (provider available, git accessible)
4. **Reviewing**: LLM API call in progress
5. **Completed**: ReviewResponse returned
6. **Failed**: Error state (invalid request, API failure, timeout)

No persistent state; each review is independent.

## Validation Summary

| Entity | Key Validation Rules |
|--------|---------------------|
| ReviewRequest | Source-type-specific field requirements, enum validations |
| ReviewResponse | Non-empty summary, positive duration |
| Finding | Valid category/severity, non-empty description/suggestion |
| ProviderConfig | Non-empty credentials, positive timeouts |
| GitDiff | Non-empty content, ≥1 file |
| DiffFile | Valid change type, ≥0 line counts |

## Next Steps

Generate MCP tool contracts in `contracts/` directory based on these entities.

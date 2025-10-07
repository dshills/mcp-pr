# Feature Specification: MCP Code Review Server

**Feature Branch**: `001-write-an-mcp`
**Created**: 2025-10-07
**Status**: Draft
**Input**: User description: "Write an MCP server for doing code reviews. Options for reviewing arbitrary code, git staged, upstaged, or specific a commit. Use the go-sdk/mcp library for the MCP. Use the OpenAI, Google, Anthropic SDKs for connecting to models for the review."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Review Arbitrary Code Snippets (Priority: P1)

A developer wants to get instant feedback on a code snippet they're working on before committing. They provide raw code text and receive a comprehensive review covering code quality, potential bugs, security issues, and best practice violations.

**Why this priority**: This is the most fundamental use case - reviewing code without requiring git integration. It delivers immediate value and can be tested completely independently of version control systems.

**Independent Test**: Can be fully tested by submitting a simple code string and verifying that a structured review response is returned with identified issues and suggestions.

**Acceptance Scenarios**:

1. **Given** a code snippet with syntax errors, **When** submitted for review, **Then** the system identifies syntax errors and suggests corrections
2. **Given** valid code with security vulnerabilities, **When** submitted for review, **Then** the system flags security issues with severity levels
3. **Given** well-written code, **When** submitted for review, **Then** the system returns a positive review with no critical issues
4. **Given** an empty code string, **When** submitted for review, **Then** the system returns a validation error

---

### User Story 2 - Review Git Staged Changes (Priority: P2)

A developer has staged changes in their git repository and wants to review them before committing. They request a review of staged changes and receive feedback on the diff between staged content and the current HEAD.

**Why this priority**: This integrates the review process directly into the commit workflow, catching issues before they enter version history. Requires git integration but is a natural workflow enhancement.

**Independent Test**: Can be tested by staging files in a test repository, requesting a staged review, and verifying that only staged changes are analyzed.

**Acceptance Scenarios**:

1. **Given** staged changes in multiple files, **When** review is requested, **Then** the system analyzes all staged diffs and returns file-by-file feedback
2. **Given** no staged changes, **When** review is requested, **Then** the system returns a message indicating nothing to review
3. **Given** staged binary files, **When** review is requested, **Then** the system skips binary files and reviews only text files

---

### User Story 3 - Review Unstaged Changes (Priority: P3)

A developer has made local modifications but hasn't staged them yet. They want preliminary feedback on work-in-progress code before deciding what to stage. The system reviews the diff between working directory and HEAD.

**Why this priority**: Enables early feedback on work in progress. Useful but less critical than staged reviews since unstaged work is typically more exploratory.

**Independent Test**: Can be tested by modifying files without staging, requesting an unstaged review, and verifying that working directory changes are analyzed.

**Acceptance Scenarios**:

1. **Given** unstaged modifications in tracked files, **When** review is requested, **Then** the system analyzes all unstaged diffs
2. **Given** new untracked files, **When** review is requested, **Then** the system optionally includes untracked files based on configuration
3. **Given** both staged and unstaged changes, **When** unstaged review is requested, **Then** the system reviews only unstaged changes

---

### User Story 4 - Review Specific Commit (Priority: P4)

A developer wants to review a historical commit to understand past decisions or audit code quality. They specify a commit SHA and receive a review of that commit's changes.

**Why this priority**: Useful for auditing and learning from past commits, but not part of the active development workflow.

**Independent Test**: Can be tested by specifying a known commit SHA and verifying that the diff for that commit is retrieved and reviewed.

**Acceptance Scenarios**:

1. **Given** a valid commit SHA, **When** review is requested, **Then** the system analyzes the commit diff and returns feedback
2. **Given** an invalid commit SHA, **When** review is requested, **Then** the system returns an error with the invalid SHA
3. **Given** a merge commit, **When** review is requested, **Then** the system reviews all changes introduced by the merge

---

### Edge Cases

- What happens when the repository has no commits (fresh init)?
- How does the system handle extremely large diffs (>10,000 lines)?
- What happens when git is not installed or not in PATH?
- How does the system handle repositories with no configured remote?
- What happens when model API calls fail or timeout?
- How does the system handle code in unsupported or uncommon languages?
- What happens when multiple API providers are unavailable?
- How does the system behave with mixed line endings (CRLF/LF)?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST accept arbitrary code text for review without requiring a git repository
- **FR-002**: System MUST retrieve and review git staged changes when requested
- **FR-003**: System MUST retrieve and review git unstaged changes when requested
- **FR-004**: System MUST retrieve and review a specific commit by SHA when requested
- **FR-005**: System MUST support multiple model providers (Anthropic, OpenAI, Google) for generating reviews
- **FR-006**: System MUST allow users to select which model provider to use for each review
- **FR-007**: System MUST return structured review results including identified issues, severity levels, and recommendations
- **FR-008**: System MUST validate git operations and return clear error messages when git commands fail
- **FR-009**: System MUST handle provider API failures gracefully with retry logic and fallback options
- **FR-010**: System MUST expose functionality through Model Context Protocol (MCP) resources and tools
- **FR-011**: System MUST categorize review findings (bugs, security, performance, style, best practices)
- **FR-012**: System MUST support configurable review depth (quick scan vs. thorough analysis)
- **FR-013**: System MUST handle code in multiple programming languages
- **FR-014**: System MUST respect diff size limits and chunk large reviews appropriately
- **FR-015**: System MUST log all review requests and responses for debugging and audit purposes

### Key Entities

- **Code Review Request**: Represents a request for code review, containing source type (arbitrary, staged, unstaged, commit), code content or git reference, provider selection, and review parameters
- **Review Response**: Contains categorized findings, each with severity level, issue description, affected code location, and suggested remediation
- **Provider Configuration**: Stores API keys, model selections, timeout settings, and retry policies for each supported provider
- **Git Context**: Repository path, current branch, commit information, and diff metadata for version-controlled reviews

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users receive review results for arbitrary code within 10 seconds for snippets under 500 lines
- **SC-002**: System successfully processes reviews using any of the three supported providers with 95% success rate
- **SC-003**: System identifies at least 80% of common security vulnerabilities in test code samples
- **SC-004**: Users can complete a full review workflow (request → receive feedback → act on feedback) in under 30 seconds
- **SC-005**: System handles git repositories with up to 100 changed files without timeout
- **SC-006**: Review responses categorize findings with 90% accuracy compared to expert manual review
- **SC-007**: System recovers from provider failures within 5 seconds by retrying or falling back to alternative providers

## Assumptions

- Users have git installed and available in PATH when using git-related features
- Users have valid API keys configured for at least one model provider
- Code being reviewed is primarily in common programming languages (Go, Python, JavaScript, etc.)
- Diff sizes are typically under 5,000 lines; extremely large diffs may require chunking
- Users understand basic git concepts (staging, commits, diffs)
- Network connectivity is available for API calls to model providers
- Review quality depends on the selected model's capabilities; responses are advisory, not authoritative

<!--
Sync Impact Report
==================
Version: 0.0.0 → 1.0.0 (Initial ratification)
Rationale: First constitution for MCP-PR project

Changes:
- Initial ratification of 7 core principles
- Added MCP Integration Standards section
- Added Development Workflow section
- Established governance rules

Modified Principles: N/A (initial version)
Added Sections:
  - Core Principles (7 principles)
  - MCP Integration Standards
  - Development Workflow
  - Governance

Removed Sections: N/A

Template Validation Status:
  ✅ .specify/templates/plan-template.md - Constitution Check section aligns
  ✅ .specify/templates/spec-template.md - Requirements sections align
  ✅ .specify/templates/tasks-template.md - Task categorization aligns
  ✅ .claude/commands/*.md - Generic guidance verified

Follow-up TODOs: None
-->

# MCP-PR Constitution

## Core Principles

### I. MCP Protocol Compliance

All features MUST adhere to the Model Context Protocol specification. Protocol
violations are unacceptable and MUST be caught in contract testing. Every
MCP-related feature MUST include protocol validation tests before implementation.

**Rationale**: MCP is the foundational protocol; non-compliance breaks
interoperability with clients and servers.

### II. Multi-Provider Support

The system MUST maintain provider-agnostic abstractions. Features MUST work
across Anthropic, OpenAI, and Google GenAI providers without provider-specific
implementations leaking into core logic. Provider adapters MUST implement a
common interface.

**Rationale**: Vendor lock-in reduces flexibility; users expect choice in LLM
providers.

### III. Go Idiomatic Design

Code MUST follow Go conventions and idioms. Use interfaces for abstraction,
prefer composition over inheritance, handle errors explicitly (no panic in
library code), and document all exported symbols. Run `gofmt`, `golint`, and
`go vet` before commits.

**Rationale**: Go community standards ensure maintainability and predictable
behavior across the ecosystem.

### IV. Test-Driven Development

Tests MUST be written before implementation. Each feature requires:
- Contract tests for MCP protocol compliance
- Integration tests for provider interactions
- Unit tests for business logic

Tests MUST fail initially, then pass after implementation (Red-Green-Refactor).

**Rationale**: TDD ensures testable design, catches regressions early, and
provides executable specifications.

### V. Observability & Debugging

All operations MUST be observable. Structured logging (JSON format) is required
for all significant events. Errors MUST include context (stack traces,
operation IDs). Performance-critical paths MUST be instrumented with metrics.

**Rationale**: Debugging distributed LLM integrations requires comprehensive
logging and tracing.

### VI. Semantic Versioning

The project MUST follow semantic versioning (MAJOR.MINOR.PATCH):
- MAJOR: Breaking changes to public APIs or MCP protocol handling
- MINOR: New features, new provider support, backward-compatible additions
- PATCH: Bug fixes, performance improvements, documentation updates

Breaking changes MUST include migration guides.

**Rationale**: Users depend on stable APIs; clear versioning prevents
unexpected breakage.

### VII. Simplicity & YAGNI

Start with the simplest implementation that solves the problem. Avoid
speculative features, premature optimization, and unnecessary abstraction.
Complex patterns (factories, repositories, complex dependency injection) MUST
be justified in the Complexity Tracking section of the implementation plan.

**Rationale**: Over-engineering increases maintenance burden and reduces
iteration speed.

## MCP Integration Standards

### Protocol Handling

- All MCP requests/responses MUST be validated against JSON-RPC 2.0
- WebSocket connections MUST implement proper reconnection logic
- Rate limiting MUST be implemented per provider specifications
- Timeouts MUST be configurable with sensible defaults

### Provider Integration

- Each provider MUST have its own adapter implementing a common interface
- API keys MUST be read from environment variables or secure configuration
- Provider-specific features MUST be documented as extensions
- Fallback behavior MUST be defined for provider failures

### Error Handling

- Network errors MUST be retried with exponential backoff
- Provider quota/rate limit errors MUST return actionable error messages
- Protocol errors MUST include the malformed request/response for debugging
- All errors MUST be logged with full context

## Development Workflow

### Feature Development

1. All features MUST start with a specification in `specs/###-feature-name/`
2. Specifications MUST NOT include implementation details (frameworks, libraries)
3. Implementation plans MUST pass Constitution Check gates
4. Tasks MUST be organized by user story (P1, P2, P3 priority order)
5. Each user story MUST be independently testable

### Code Quality Gates

- All tests MUST pass before merge
- Code coverage MUST be ≥80% for new code
- `gofmt`, `go vet`, and `golangci-lint` MUST run clean
- All exported functions MUST have godoc comments
- Breaking changes MUST be documented in CHANGELOG.md

### Testing Requirements

- Contract tests MUST verify MCP protocol compliance
- Integration tests MUST cover all provider interactions
- Mock external dependencies in unit tests
- Performance tests for operations handling >100 requests/second

## Governance

This constitution supersedes all other development practices and guidelines.
Amendments require:
1. Documented rationale for the change
2. Impact assessment on existing code
3. Updated Sync Impact Report
4. Version bump following semantic versioning rules

All pull requests and code reviews MUST verify compliance with these
principles. Violations MUST be flagged and justified in the Complexity Tracking
section of the implementation plan, or corrected before merge.

For runtime development guidance, refer to `CLAUDE.md` (for Claude Code) or
equivalent agent-specific files. These guidance files MUST remain consistent
with this constitution.

**Version**: 1.0.0 | **Ratified**: 2025-10-07 | **Last Amended**: 2025-10-07

# Implementation Plan: MCP Code Review Server

**Branch**: `001-write-an-mcp` | **Date**: 2025-10-07 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-write-an-mcp/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Build an MCP server that provides code review capabilities through the Model Context Protocol. The server will accept code from four sources (arbitrary text, git staged changes, unstaged changes, or specific commits) and generate structured reviews using Anthropic, OpenAI, or Google GenAI models. Reviews categorize findings by type (bugs, security, performance, style) with severity levels and remediation suggestions.

## Technical Context

**Language/Version**: Go 1.25.1
**Primary Dependencies**:
- github.com/modelcontextprotocol/go-sdk v1.0.0 (MCP server implementation)
- github.com/anthropics/anthropic-sdk-go v1.13.0 (Anthropic Claude API)
- github.com/openai/openai-go/v3 v3.2.0 (OpenAI API)
- google.golang.org/genai v1.28.0 (Google Gemini API)
- Standard library: os/exec for git operations, encoding/json for structured data

**Storage**: N/A (stateless server; configuration from environment variables)
**Testing**: Go standard testing package, table-driven tests, mock provider interfaces
**Target Platform**: Cross-platform (Linux, macOS, Windows) - MCP server runs as subprocess or stdio transport
**Project Type**: Single project (MCP server binary)
**Performance Goals**: <10s review latency for 500-line code snippets; handle diffs up to 5,000 lines
**Constraints**: Network-dependent (requires LLM API access); git must be in PATH for repository operations
**Scale/Scope**: Single-user tool (MCP servers run per-client); stateless request/response model

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Compliance Notes |
|-----------|--------|------------------|
| I. MCP Protocol Compliance | ✅ PASS | Feature explicitly uses go-sdk/mcp library; will include MCP contract tests |
| II. Multi-Provider Support | ✅ PASS | Design requires 3 providers (Anthropic, OpenAI, Google) with common interface |
| III. Go Idiomatic Design | ✅ PASS | Using Go 1.25.1 with standard project structure; will follow gofmt/golint/go vet |
| IV. Test-Driven Development | ✅ PASS | Plan includes contract, integration, and unit test phases before implementation |
| V. Observability & Debugging | ✅ PASS | FR-015 requires structured logging; will implement JSON logging for all operations |
| VI. Semantic Versioning | ✅ PASS | Initial version will be v0.1.0; plan includes CHANGELOG.md |
| VII. Simplicity & YAGNI | ✅ PASS | Simple adapter pattern for providers; no complex dependency injection |

**Overall**: ✅ PASS - All constitution principles satisfied. No violations requiring justification.

## Project Structure

### Documentation (this feature)

```
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```
cmd/
└── mcp-code-review/
    └── main.go                    # MCP server entry point

internal/
├── mcp/
│   ├── server.go                  # MCP server implementation
│   ├── tools.go                   # MCP tool handlers (review_code, review_staged, etc.)
│   └── resources.go               # MCP resource handlers (if needed)
├── providers/
│   ├── provider.go                # Provider interface
│   ├── anthropic.go               # Anthropic adapter
│   ├── openai.go                  # OpenAI adapter
│   └── google.go                  # Google GenAI adapter
├── review/
│   ├── engine.go                  # Core review orchestration
│   ├── request.go                 # Review request models
│   └── response.go                # Review response models
├── git/
│   ├── client.go                  # Git command wrapper
│   ├── diff.go                    # Diff parsing and formatting
│   └── repo.go                    # Repository operations
├── config/
│   └── config.go                  # Configuration loading (env vars)
└── logging/
    └── logger.go                  # Structured logging

tests/
├── contract/
│   ├── mcp_protocol_test.go       # MCP JSON-RPC compliance tests
│   └── provider_interface_test.go # Provider contract tests
├── integration/
│   ├── anthropic_test.go          # Live Anthropic integration
│   ├── openai_test.go             # Live OpenAI integration
│   ├── google_test.go             # Live Google integration
│   └── git_test.go                # Git operations in test repos
└── unit/
    ├── review_test.go             # Review engine unit tests
    ├── git_test.go                # Git client unit tests (mocked exec)
    └── mcp_test.go                # MCP handlers unit tests

go.mod
go.sum
README.md
CHANGELOG.md
```

**Structure Decision**: Single project structure selected. This is an MCP server binary with clear separation of concerns: MCP protocol layer (`internal/mcp/`), provider adapters (`internal/providers/`), core review logic (`internal/review/`), git operations (`internal/git/`), and configuration/logging utilities. The `cmd/` directory contains the server entry point, while `internal/` prevents external package dependencies.

## Complexity Tracking

*Fill ONLY if Constitution Check has violations that must be justified*

No violations. Constitution Check passed all gates.

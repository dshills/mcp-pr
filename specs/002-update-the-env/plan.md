# Implementation Plan: Project-Specific Environment Variables

**Branch**: `002-update-the-env` | **Date**: 2025-10-29 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-update-the-env/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Rename generic `MCP_*` environment variables to project-specific `MCP_PR_*` names to prevent namespace collisions when multiple MCP servers are installed. The feature preserves standard LLM API key conventions (`ANTHROPIC_API_KEY`, `OPENAI_API_KEY`, `GOOGLE_API_KEY`) and provider-specific timeouts while adding backward compatibility support during the transition period.

## Technical Context

**Language/Version**: Go 1.25.1
**Primary Dependencies**: Standard library (os, time, fmt packages for environment variable handling)
**Storage**: N/A (configuration is loaded from environment variables at startup)
**Testing**: Go standard testing (testing package), contract tests for configuration loading
**Target Platform**: Cross-platform (Linux, macOS, Windows) - runs as MCP server subprocess
**Project Type**: Single project (MCP server binary)
**Performance Goals**: Configuration loading must complete in <10ms
**Constraints**: Must maintain backward compatibility for at least one major version
**Scale/Scope**: 4 environment variables to rename, ~20 files to update (code + documentation)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### I. MCP Protocol Compliance ✅
**Status**: PASS - Configuration changes do not affect MCP protocol compliance. Environment variables are internal configuration, not protocol-facing.

### II. Multi-Provider Support ✅
**Status**: PASS - Preserving standard API key names (`ANTHROPIC_API_KEY`, `OPENAI_API_KEY`, `GOOGLE_API_KEY`) and provider-specific timeouts maintains provider-agnostic design. Only general server configuration variables are being renamed.

### III. Go Idiomatic Design ✅
**Status**: PASS - Using standard library `os.Getenv()` pattern with fallback defaults is idiomatic Go. The `getEnv()` helper function follows Go conventions.

### IV. Test-Driven Development ✅
**Status**: PASS - Will write tests first for:
- Configuration loading with new variable names
- Backward compatibility (both old and new variables)
- Precedence (new variables override old when both set)
- Default value fallback behavior

### V. Observability & Debugging ✅
**Status**: PASS - Configuration values (excluding API keys) are already logged at startup. Will add deprecation warnings when old variable names are detected.

### VI. Semantic Versioning ✅
**Status**: PASS - This is a MINOR version bump (backward-compatible feature addition). Adds new configuration variables while maintaining support for old ones. No breaking changes during transition period.

### VII. Simplicity & YAGNI ✅
**Status**: PASS - Straightforward rename with backward compatibility check. No new abstractions or complex patterns needed.

### Additional Gates

**Code Quality Gates** ✅
- Tests pass: Will add unit tests before implementation
- Coverage ≥80%: Targeting 100% for config package changes
- gofmt/go vet/golangci-lint: Will run before commit
- Godoc: Will update for modified functions
- CHANGELOG: Will document as backward-compatible enhancement

**Testing Requirements** ✅
- Contract tests: N/A (no MCP protocol changes)
- Integration tests: Will update existing tests to use new variable names
- Unit tests: Will add tests for backward compatibility logic
- Performance tests: N/A (config loading is initialization-only, not high-throughput)

**All gates PASS** - No violations to justify in Complexity Tracking.

## Project Structure

### Documentation (this feature)

```
specs/002-update-the-env/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```
mcp-pr/
├── cmd/
│   └── mcp-code-review/
│       └── main.go                    # Entry point (may need updates if logging config)
├── internal/
│   ├── config/
│   │   └── config.go                  # PRIMARY: Load() function to update
│   ├── logging/
│   │   └── logger.go                  # May need updates for deprecation warnings
│   └── [other packages unchanged]
├── tests/
│   ├── integration/
│   │   ├── anthropic_test.go          # Update to use new variable names
│   │   ├── openai_test.go             # Update to use new variable names
│   │   ├── google_test.go             # Update to use new variable names
│   │   └── helpers.go                 # Update setup/teardown if setting env vars
│   ├── unit/
│   │   └── [config tests to add]     # New tests for backward compatibility
│   └── contract/
│       └── [unchanged]
├── .github/
│   └── workflows/
│       └── ci.yml                      # Update environment variable setup
├── README.md                           # Update configuration section
├── CONTRIBUTING.md                     # Update development setup if referenced
└── CHANGELOG.md                        # Document changes
```

**Structure Decision**: Single project structure. All changes are localized to:
1. Configuration loading (`internal/config/config.go`)
2. Test files that set environment variables
3. Documentation files (README, CONTRIBUTING, CI workflows)
4. No new directories or files needed beyond test files

## Complexity Tracking

*Fill ONLY if Constitution Check has violations that must be justified*

N/A - All constitution gates pass. No complexity violations to justify.

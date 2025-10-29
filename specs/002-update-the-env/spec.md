# Feature Specification: Project-Specific Environment Variables

**Feature Branch**: `002-update-the-env`
**Created**: 2025-10-29
**Status**: Draft
**Input**: User description: "update the env variables to make them specific to the project. for example MCP_LOG_LEVEL should be MCP_PR_LOG_LEVEL. The exception is the LLM API Keys"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Namespace Environment Variables (Priority: P1)

When a developer installs the mcp-pr server alongside other MCP servers, they need project-specific environment variables to avoid configuration conflicts. Currently, generic names like `MCP_LOG_LEVEL` and `MCP_DEFAULT_PROVIDER` could collide with other MCP server implementations.

**Why this priority**: This is the core requirement that prevents namespace pollution and allows multiple MCP servers to coexist with independent configurations. Without this, users cannot reliably configure multiple MCP servers.

**Independent Test**: Can be fully tested by setting `MCP_PR_LOG_LEVEL=debug` while other MCP servers use `MCP_LOG_LEVEL=error`, and verifying that mcp-pr uses debug-level logging independently.

**Acceptance Scenarios**:

1. **Given** a developer has multiple MCP servers installed, **When** they set `MCP_PR_LOG_LEVEL=debug` and `MCP_LOG_LEVEL=error`, **Then** mcp-pr server uses debug logging while other servers use error logging
2. **Given** a developer configures `MCP_PR_DEFAULT_PROVIDER=openai`, **When** they start mcp-pr, **Then** it defaults to OpenAI provider without affecting other MCP servers' provider selection
3. **Given** a developer sets `MCP_PR_REVIEW_TIMEOUT=180s`, **When** they perform a code review, **Then** the timeout is applied without affecting other MCP servers' timeout settings

---

### User Story 2 - Maintain LLM API Key Compatibility (Priority: P1)

Developers should continue using standard LLM provider API key environment variables (`ANTHROPIC_API_KEY`, `OPENAI_API_KEY`, `GOOGLE_API_KEY`) because these are industry-standard conventions shared across tools and SDKs.

**Why this priority**: API keys are the authentication mechanism and must remain compatible with existing developer workflows and documentation. Breaking this convention would create unnecessary friction and confusion.

**Independent Test**: Can be tested by configuring standard API keys in environment and verifying all three providers authenticate successfully with existing key names.

**Acceptance Scenarios**:

1. **Given** a developer has `ANTHROPIC_API_KEY` set in their environment, **When** they start mcp-pr with `provider=anthropic`, **Then** authentication succeeds using the standard key name
2. **Given** a developer has all three provider keys set, **When** they switch between providers, **Then** each provider authenticates using its standard key name without requiring project-specific key variables
3. **Given** a developer follows standard Claude/OpenAI/Google documentation for setting API keys, **When** they configure mcp-pr, **Then** their existing key environment variables work without modification

---

### User Story 3 - Update Documentation (Priority: P2)

Developers need clear documentation of the new environment variable names to configure the server correctly after upgrading.

**Why this priority**: Documentation is essential for adoption but can be completed after the code changes. It doesn't block functionality but is required before release.

**Independent Test**: Can be tested by following the README instructions to configure a new installation and verifying all environment variables work as documented.

**Acceptance Scenarios**:

1. **Given** a developer reads the README, **When** they configure the server for the first time, **Then** they understand all available environment variables and their defaults
2. **Given** an existing user upgrades to the new version, **When** they read the migration guide, **Then** they know exactly which environment variables to rename in their configuration
3. **Given** a developer troubleshoots configuration issues, **When** they check the documentation, **Then** they can find examples showing both old and new variable names with clear migration instructions

---

### Edge Cases

- What happens when developers use old environment variable names (`MCP_LOG_LEVEL`) after upgrade?
- How does the system handle when both old and new variable names are set (e.g., `MCP_LOG_LEVEL` and `MCP_PR_LOG_LEVEL`)?
- What happens when a developer has multiple MCP servers that use the same `MCP_` prefix?
- How does the server behave when provider-specific timeouts are set but review timeout is not?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST rename `MCP_LOG_LEVEL` to `MCP_PR_LOG_LEVEL` for log level configuration
- **FR-002**: System MUST rename `MCP_DEFAULT_PROVIDER` to `MCP_PR_DEFAULT_PROVIDER` for provider selection
- **FR-003**: System MUST rename `MCP_REVIEW_TIMEOUT` to `MCP_PR_REVIEW_TIMEOUT` for timeout configuration
- **FR-004**: System MUST rename `MCP_MAX_DIFF_SIZE` to `MCP_PR_MAX_DIFF_SIZE` for diff size limits
- **FR-005**: System MUST preserve standard LLM provider API key names: `ANTHROPIC_API_KEY`, `OPENAI_API_KEY`, `GOOGLE_API_KEY`
- **FR-006**: System MUST preserve provider-specific timeout names: `ANTHROPIC_TIMEOUT`, `OPENAI_TIMEOUT`, `GOOGLE_TIMEOUT`
- **FR-007**: System MUST update all code references to use the new environment variable names
- **FR-008**: System MUST update README.md with the new environment variable names
- **FR-009**: System MUST update configuration examples (Claude Desktop config) with new variable names
- **FR-010**: System MUST update CONTRIBUTING.md with new environment variable names if referenced
- **FR-011**: System MUST update test files to use new environment variable names
- **FR-012**: System MUST maintain backward compatibility during transition by checking both old and new variable names, preferring new names when both are set

### Key Entities

- **Configuration**: Collection of environment variables including log level, provider selection, timeouts, and size limits
- **Environment Variable**: Key-value pair loaded from system environment with specific naming convention (`MCP_PR_*` prefix)
- **API Key**: Authentication credential for LLM providers using standard industry naming conventions

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Developers can configure mcp-pr server with project-specific variables without conflicts when running multiple MCP servers
- **SC-002**: Existing installations continue working during transition period with backward compatibility support
- **SC-003**: All documentation accurately reflects new environment variable names within one day of code changes
- **SC-004**: 100% of code references to generic `MCP_*` variables (excluding API keys and provider timeouts) are updated to `MCP_PR_*` prefix
- **SC-005**: Developers can successfully authenticate with all three LLM providers using standard API key environment variables
- **SC-006**: Configuration loading completes without errors when using new variable names

## Assumptions *(if applicable)*

- Developers will need a transition period where both old and new variable names are supported
- The project name "mcp-pr" is the appropriate namespace prefix
- Provider-specific timeouts (`ANTHROPIC_TIMEOUT`, etc.) do not need project prefixing because they are already provider-namespaced
- Standard API key names are widely adopted across the ecosystem and should not be changed
- Users are expected to manage environment variables through their shell, CI/CD, or MCP client configuration files
- Backward compatibility will be maintained for at least one major version before removing support for old variable names

## Out of Scope *(if applicable)*

- Changing the naming convention for API keys (they remain as industry-standard names)
- Creating a configuration file format (remaining with environment variables only)
- Migrating existing user configurations automatically
- Adding new environment variables beyond renaming existing ones
- Changing provider-specific timeout variable names

## Dependencies *(if applicable)*

- Requires Go code updates in `internal/config/config.go`
- Requires documentation updates in `README.md`, `CONTRIBUTING.md`, and configuration examples
- Requires test file updates to reflect new variable names
- May require updates to CI/CD workflows in `.github/workflows/ci.yml`

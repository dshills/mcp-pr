# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial MCP server implementation
- Support for Anthropic Claude, OpenAI GPT, and Google Gemini providers
- `review_code` tool for arbitrary code review
- `review_staged` tool for git staged changes
- `review_unstaged` tool for git unstaged changes
- `review_commit` tool for specific commit review
- Structured JSON logging
- Multi-provider adapter pattern
- TDD test suite (contract, integration, unit tests)
- Project-specific environment variable names (`MCP_PR_*` prefix)
- Backward compatibility support for old environment variable names
- Deprecation warnings for old variable names

### Changed
- **MINOR**: Renamed environment variables to project-specific names:
  - `MCP_LOG_LEVEL` → `MCP_PR_LOG_LEVEL` (old name deprecated, will be removed in v1.0.0)
  - `MCP_DEFAULT_PROVIDER` → `MCP_PR_DEFAULT_PROVIDER` (old name deprecated, will be removed in v1.0.0)
  - `MCP_REVIEW_TIMEOUT` → `MCP_PR_REVIEW_TIMEOUT` (old name deprecated, will be removed in v1.0.0)
  - `MCP_MAX_DIFF_SIZE` → `MCP_PR_MAX_DIFF_SIZE` (old name deprecated, will be removed in v1.0.0)
- Old variable names still work during transition period (with deprecation warnings)
- New variable names take precedence when both old and new are set

### Migration Guide
- **API keys unchanged**: `ANTHROPIC_API_KEY`, `OPENAI_API_KEY`, `GOOGLE_API_KEY` remain the same
- **Provider timeouts unchanged**: `ANTHROPIC_TIMEOUT`, `OPENAI_TIMEOUT`, `GOOGLE_TIMEOUT` remain the same
- **Action required**: Update your configuration to use `MCP_PR_*` variable names
- **Timeline**: Old names will be removed in v1.0.0 release

## [0.1.0] - TBD

### Added
- MVP release with arbitrary code review capability

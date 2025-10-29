# Data Model: Project-Specific Environment Variables

**Feature**: 002-update-the-env
**Date**: 2025-10-29
**Status**: Complete

## Overview

This feature involves configuration data loaded from environment variables. There is no persistent storage or complex entity relationships. The data model is simple: environment variables are read at startup and mapped to a configuration struct.

## Entities

### Configuration

**Description**: Aggregates all server configuration settings loaded from environment variables at application startup.

**Attributes**:
- `AnthropicAPIKey` (string): API key for Anthropic provider, loaded from `ANTHROPIC_API_KEY`
- `OpenAIAPIKey` (string): API key for OpenAI provider, loaded from `OPENAI_API_KEY`
- `GoogleAPIKey` (string): API key for Google provider, loaded from `GOOGLE_API_KEY`
- `LogLevel` (string): Logging verbosity level (debug|info|warn|error), loaded from `MCP_PR_LOG_LEVEL` with fallback to `MCP_LOG_LEVEL`
- `DefaultProvider` (string): Default LLM provider selection (anthropic|openai|google), loaded from `MCP_PR_DEFAULT_PROVIDER` with fallback to `MCP_DEFAULT_PROVIDER`
- `ReviewTimeout` (duration): Maximum time for review operations, loaded from `MCP_PR_REVIEW_TIMEOUT` with fallback to `MCP_REVIEW_TIMEOUT`
- `MaxDiffSize` (integer): Maximum diff size in bytes, loaded from `MCP_PR_MAX_DIFF_SIZE` with fallback to `MCP_MAX_DIFF_SIZE`
- `AnthropicTimeout` (duration): Anthropic-specific API timeout, loaded from `ANTHROPIC_TIMEOUT`
- `OpenAITimeout` (duration): OpenAI-specific API timeout, loaded from `OPENAI_TIMEOUT`
- `GoogleTimeout` (duration): Google-specific API timeout, loaded from `GOOGLE_TIMEOUT`

**Validation Rules**:
- At least one API key (Anthropic, OpenAI, or Google) must be present
- LogLevel must be one of: debug, info, warn, error (default: info)
- DefaultProvider must be one of: anthropic, openai, google (default: anthropic)
- Timeouts must be valid duration strings (e.g., "90s", "2m") and greater than 0
- MaxDiffSize must be a positive integer

**Default Values**:
- LogLevel: "info"
- DefaultProvider: "anthropic"
- ReviewTimeout: 120 seconds
- MaxDiffSize: 10000 bytes
- AnthropicTimeout: 90 seconds
- OpenAITimeout: 90 seconds
- GoogleTimeout: 90 seconds

### Environment Variable

**Description**: Represents a single environment variable with its current and deprecated names.

**Attributes**:
- `NewName` (string): Current, project-specific variable name (e.g., "MCP_PR_LOG_LEVEL")
- `OldName` (string): Deprecated, generic variable name (e.g., "MCP_LOG_LEVEL")
- `DefaultValue` (string): Fallback value if neither new nor old variable is set
- `IsDeprecated` (boolean): True if old variable name is being used

**Relationships**:
- Maps to one field in the Configuration entity
- No relationships with other entities (standalone configuration)

**Lifecycle**:
1. **Load**: Read at application startup via `config.Load()`
2. **Validate**: Check for required values and format constraints
3. **Immutable**: Configuration does not change during application runtime
4. **Log Warnings**: If old variable names are detected, log deprecation warnings

## State Transitions

### Configuration Loading State Machine

```
[Application Start]
        ↓
[Load Environment Variables]
        ↓
[Check New Variable Name] → Found? → [Use New Value] → [No Warning]
        ↓
       Not Found
        ↓
[Check Old Variable Name] → Found? → [Use Old Value] → [Log Deprecation Warning]
        ↓
       Not Found
        ↓
[Use Default Value] → [No Warning]
        ↓
[Validate Configuration]
        ↓
   Valid? → [Application Ready]
        ↓
      Invalid
        ↓
[Return Error] → [Application Fails to Start]
```

### Variable Precedence

When loading each configuration value:
1. **Priority 1**: New project-specific variable (e.g., `MCP_PR_LOG_LEVEL`)
2. **Priority 2**: Old generic variable (e.g., `MCP_LOG_LEVEL`) with deprecation warning
3. **Priority 3**: Default value

## Data Flow

```
Environment Variables
        ↓
[os.Getenv() calls]
        ↓
[getEnvWithFallback() helper]
        ↓
[Configuration Struct]
        ↓
[Validation]
        ↓
[Application Components]
```

**Key Points**:
- Environment variables are read-only inputs
- No data persistence or state storage
- Configuration is loaded once at startup
- Deprecation warnings are side effects during loading

## Validation Rules by Field

| Field | New Env Var | Old Env Var (Deprecated) | Valid Values | Default |
|-------|-------------|--------------------------|--------------|---------|
| LogLevel | MCP_PR_LOG_LEVEL | MCP_LOG_LEVEL | debug\|info\|warn\|error | info |
| DefaultProvider | MCP_PR_DEFAULT_PROVIDER | MCP_DEFAULT_PROVIDER | anthropic\|openai\|google | anthropic |
| ReviewTimeout | MCP_PR_REVIEW_TIMEOUT | MCP_REVIEW_TIMEOUT | Valid duration > 0 | 120s |
| MaxDiffSize | MCP_PR_MAX_DIFF_SIZE | MCP_MAX_DIFF_SIZE | Integer > 0 | 10000 |
| AnthropicTimeout | ANTHROPIC_TIMEOUT | N/A (not renamed) | Valid duration > 0 | 90s |
| OpenAITimeout | OPENAI_TIMEOUT | N/A (not renamed) | Valid duration > 0 | 90s |
| GoogleTimeout | GOOGLE_TIMEOUT | N/A (not renamed) | Valid duration > 0 | 90s |
| AnthropicAPIKey | ANTHROPIC_API_KEY | N/A (not renamed) | Non-empty string | "" |
| OpenAIAPIKey | OPENAI_API_KEY | N/A (not renamed) | Non-empty string | "" |
| GoogleAPIKey | GOOGLE_API_KEY | N/A (not renamed) | Non-empty string | "" |

## Error Handling

**Validation Errors**:
- **No API Keys**: If all three provider API keys are empty, return error: "at least one provider API key must be configured"
- **Invalid Duration**: If timeout values cannot be parsed, use default value (graceful degradation)
- **Invalid Integer**: If MaxDiffSize cannot be parsed, use default value (graceful degradation)
- **Invalid Enum**: If LogLevel or DefaultProvider has invalid value, use default value (graceful degradation)

**Deprecation Warnings** (not errors):
- Logged at WARN level when old variable names are used
- Includes migration guidance
- Does not prevent application startup

## No Persistent Storage

This feature does not involve:
- Databases
- File storage
- Caching
- Session state
- User data persistence

All data is sourced from environment variables and held in memory during application runtime.

## Summary

The data model is minimal:
- **One entity**: Configuration struct
- **No relationships**: Flat structure of independent settings
- **No persistence**: Environment variables only
- **No state transitions**: Configuration is immutable after loading
- **Simple validation**: Type checking and required field validation

This simplicity aligns with the Constitution principle VII (Simplicity & YAGNI) - no unnecessary abstraction or complexity.

# Configuration Contract: Environment Variables

**Feature**: 002-update-the-env
**Date**: 2025-10-29
**Type**: Internal Configuration Contract

## Overview

This document defines the contract for environment variable configuration in mcp-pr. It specifies which variables are supported, their formats, validation rules, and backward compatibility guarantees.

## Scope

This is an **internal contract** - it defines how the application loads and validates configuration from the environment. This is not an external API contract exposed to clients.

## Contract Version

**Version**: 1.0.0 (introduces project-specific naming)
**Backward Compatibility**: Supported until v1.0.0 (old variable names work with deprecation warnings)

---

## Environment Variables

### Required Variables (at least one)

At least one LLM provider API key must be set. The application will fail to start if none are provided.

#### ANTHROPIC_API_KEY

- **Type**: String
- **Format**: `sk-ant-...`
- **Required**: One of Anthropic, OpenAI, or Google API key must be set
- **Default**: None (must be explicitly set)
- **Example**: `export ANTHROPIC_API_KEY="sk-ant-api03-..."`
- **Status**: Stable (not renamed, industry standard)

#### OPENAI_API_KEY

- **Type**: String
- **Format**: `sk-...`
- **Required**: One of Anthropic, OpenAI, or Google API key must be set
- **Default**: None (must be explicitly set)
- **Example**: `export OPENAI_API_KEY="sk-proj-..."`
- **Status**: Stable (not renamed, industry standard)

#### GOOGLE_API_KEY

- **Type**: String
- **Format**: Google API key format
- **Required**: One of Anthropic, OpenAI, or Google API key must be set
- **Default**: None (must be explicitly set)
- **Example**: `export GOOGLE_API_KEY="AIza..."`
- **Status**: Stable (not renamed, industry standard)

---

### Optional Variables (Project-Specific)

These variables control mcp-pr server behavior. All have sensible defaults.

#### MCP_PR_LOG_LEVEL

- **Type**: String (enum)
- **Valid Values**: `debug`, `info`, `warn`, `error`
- **Default**: `info`
- **Example**: `export MCP_PR_LOG_LEVEL=debug`
- **Status**: Current (v1.0.0+)
- **Deprecated Alternative**: `MCP_LOG_LEVEL` (will be removed in v1.0.0)
- **Description**: Controls logging verbosity for the mcp-pr server

**Backward Compatibility**:
```bash
# New (preferred)
export MCP_PR_LOG_LEVEL=debug

# Old (deprecated, works with warning)
export MCP_LOG_LEVEL=debug

# Precedence: If both are set, MCP_PR_LOG_LEVEL takes precedence
```

#### MCP_PR_DEFAULT_PROVIDER

- **Type**: String (enum)
- **Valid Values**: `anthropic`, `openai`, `google`
- **Default**: `anthropic`
- **Example**: `export MCP_PR_DEFAULT_PROVIDER=openai`
- **Status**: Current (v1.0.0+)
- **Deprecated Alternative**: `MCP_DEFAULT_PROVIDER` (will be removed in v1.0.0)
- **Description**: Selects which LLM provider to use by default when multiple API keys are configured

**Backward Compatibility**:
```bash
# New (preferred)
export MCP_PR_DEFAULT_PROVIDER=anthropic

# Old (deprecated, works with warning)
export MCP_DEFAULT_PROVIDER=anthropic
```

#### MCP_PR_REVIEW_TIMEOUT

- **Type**: Duration string
- **Format**: Go duration format (e.g., "120s", "2m", "1m30s")
- **Default**: `120s`
- **Example**: `export MCP_PR_REVIEW_TIMEOUT=180s`
- **Status**: Current (v1.0.0+)
- **Deprecated Alternative**: `MCP_REVIEW_TIMEOUT` (will be removed in v1.0.0)
- **Description**: Maximum time allowed for a review operation before timing out

**Backward Compatibility**:
```bash
# New (preferred)
export MCP_PR_REVIEW_TIMEOUT=120s

# Old (deprecated, works with warning)
export MCP_REVIEW_TIMEOUT=120s
```

#### MCP_PR_MAX_DIFF_SIZE

- **Type**: Integer
- **Format**: Positive integer (bytes)
- **Default**: `10000`
- **Example**: `export MCP_PR_MAX_DIFF_SIZE=20000`
- **Status**: Current (v1.0.0+)
- **Deprecated Alternative**: `MCP_MAX_DIFF_SIZE` (will be removed in v1.0.0)
- **Description**: Maximum diff size in bytes that the server will process

**Backward Compatibility**:
```bash
# New (preferred)
export MCP_PR_MAX_DIFF_SIZE=10000

# Old (deprecated, works with warning)
export MCP_MAX_DIFF_SIZE=10000
```

---

### Optional Variables (Provider-Specific Timeouts)

#### ANTHROPIC_TIMEOUT

- **Type**: Duration string
- **Default**: `90s`
- **Example**: `export ANTHROPIC_TIMEOUT=120s`
- **Status**: Stable (not renamed, provider-namespaced)
- **Description**: Timeout for Anthropic API calls

#### OPENAI_TIMEOUT

- **Type**: Duration string
- **Default**: `90s`
- **Example**: `export OPENAI_TIMEOUT=120s`
- **Status**: Stable (not renamed, provider-namespaced)
- **Description**: Timeout for OpenAI API calls

#### GOOGLE_TIMEOUT

- **Type**: Duration string
- **Default**: `90s`
- **Example**: `export GOOGLE_TIMEOUT=120s`
- **Status**: Stable (not renamed, provider-namespaced)
- **Description**: Timeout for Google GenAI API calls

---

## Validation Rules

### Startup Validation

The following validation occurs during `config.Load()` at application startup:

1. **At least one API key**: At least one of `ANTHROPIC_API_KEY`, `OPENAI_API_KEY`, or `GOOGLE_API_KEY` must be set
   - **Error if violated**: `"at least one provider API key must be configured"`
   - **Application behavior**: Fails to start

2. **Log level enum**: If set, `MCP_PR_LOG_LEVEL` must be one of: debug, info, warn, error
   - **Error if violated**: None (falls back to default "info")
   - **Application behavior**: Continues with default value

3. **Provider enum**: If set, `MCP_PR_DEFAULT_PROVIDER` must be one of: anthropic, openai, google
   - **Error if violated**: None (falls back to default "anthropic")
   - **Application behavior**: Continues with default value

4. **Duration format**: All timeout values must be valid Go duration strings
   - **Error if violated**: None (falls back to default values)
   - **Application behavior**: Continues with default value

5. **Positive integers**: `MCP_PR_MAX_DIFF_SIZE` must be a positive integer
   - **Error if violated**: None (falls back to default 10000)
   - **Application behavior**: Continues with default value

### Runtime Validation

Configuration is loaded once at startup and is immutable during runtime. No runtime validation occurs.

---

## Backward Compatibility Guarantees

### Transition Period (Current Version → v1.0.0)

**Old Variable Names Supported**:
- `MCP_LOG_LEVEL` → Use `MCP_PR_LOG_LEVEL` instead
- `MCP_DEFAULT_PROVIDER` → Use `MCP_PR_DEFAULT_PROVIDER` instead
- `MCP_REVIEW_TIMEOUT` → Use `MCP_PR_REVIEW_TIMEOUT` instead
- `MCP_MAX_DIFF_SIZE` → Use `MCP_PR_MAX_DIFF_SIZE` instead

**Behavior**:
1. Old variable names still work
2. Deprecation warning logged at startup (WARN level)
3. New variable names take precedence if both are set
4. Application functionality is unaffected

**Deprecation Warning Format**:
```
WARN: Environment variable "MCP_LOG_LEVEL" is deprecated and will be removed in v1.0.0.
Please use "MCP_PR_LOG_LEVEL" instead.
```

### Breaking Changes in v1.0.0

**Support Removed**:
- `MCP_LOG_LEVEL` (must use `MCP_PR_LOG_LEVEL`)
- `MCP_DEFAULT_PROVIDER` (must use `MCP_PR_DEFAULT_PROVIDER`)
- `MCP_REVIEW_TIMEOUT` (must use `MCP_PR_REVIEW_TIMEOUT`)
- `MCP_MAX_DIFF_SIZE` (must use `MCP_PR_MAX_DIFF_SIZE`)

**Not Affected** (stable across all versions):
- `ANTHROPIC_API_KEY`
- `OPENAI_API_KEY`
- `GOOGLE_API_KEY`
- `ANTHROPIC_TIMEOUT`
- `OPENAI_TIMEOUT`
- `GOOGLE_TIMEOUT`

---

## Testing Contract

### Unit Tests

Unit tests must verify:
1. ✅ New variable names load correctly
2. ✅ Old variable names load correctly (with deprecation warning)
3. ✅ New variables take precedence when both are set
4. ✅ Default values are used when neither is set
5. ✅ Validation errors are handled gracefully
6. ✅ At least one API key is required

### Integration Tests

Integration tests must verify:
1. ✅ All three providers work with standard API key names
2. ✅ Configuration loads correctly in real startup scenario
3. ✅ Deprecation warnings appear in logs when old variables are used

---

## Configuration Examples

### Minimal Configuration (New Variables)

```bash
export ANTHROPIC_API_KEY="sk-ant-..."
export MCP_PR_LOG_LEVEL=info
export MCP_PR_DEFAULT_PROVIDER=anthropic
```

### Multi-Provider Configuration (New Variables)

```bash
export ANTHROPIC_API_KEY="sk-ant-..."
export OPENAI_API_KEY="sk-..."
export GOOGLE_API_KEY="AIza..."
export MCP_PR_DEFAULT_PROVIDER=anthropic
export MCP_PR_LOG_LEVEL=debug
export MCP_PR_REVIEW_TIMEOUT=180s
export MCP_PR_MAX_DIFF_SIZE=20000
```

### Claude Desktop Configuration (New Variables)

```json
{
  "mcpServers": {
    "code-review": {
      "command": "/path/to/mcp-code-review",
      "env": {
        "ANTHROPIC_API_KEY": "sk-ant-...",
        "MCP_PR_LOG_LEVEL": "info",
        "MCP_PR_DEFAULT_PROVIDER": "anthropic",
        "MCP_PR_REVIEW_TIMEOUT": "120s"
      }
    }
  }
}
```

### Backward Compatible Configuration (Deprecated)

```bash
# This works but logs deprecation warnings
export ANTHROPIC_API_KEY="sk-ant-..."
export MCP_LOG_LEVEL=info              # DEPRECATED: Use MCP_PR_LOG_LEVEL
export MCP_DEFAULT_PROVIDER=anthropic  # DEPRECATED: Use MCP_PR_DEFAULT_PROVIDER
```

---

## Summary Table

| Variable | Type | Default | Status | Notes |
|----------|------|---------|--------|-------|
| ANTHROPIC_API_KEY | string | none | Stable | Industry standard |
| OPENAI_API_KEY | string | none | Stable | Industry standard |
| GOOGLE_API_KEY | string | none | Stable | Industry standard |
| MCP_PR_LOG_LEVEL | enum | info | Current | Replaces MCP_LOG_LEVEL |
| MCP_PR_DEFAULT_PROVIDER | enum | anthropic | Current | Replaces MCP_DEFAULT_PROVIDER |
| MCP_PR_REVIEW_TIMEOUT | duration | 120s | Current | Replaces MCP_REVIEW_TIMEOUT |
| MCP_PR_MAX_DIFF_SIZE | integer | 10000 | Current | Replaces MCP_MAX_DIFF_SIZE |
| ANTHROPIC_TIMEOUT | duration | 90s | Stable | Provider-specific |
| OPENAI_TIMEOUT | duration | 90s | Stable | Provider-specific |
| GOOGLE_TIMEOUT | duration | 90s | Stable | Provider-specific |
| MCP_LOG_LEVEL | enum | N/A | Deprecated | Remove in v1.0.0 |
| MCP_DEFAULT_PROVIDER | enum | N/A | Deprecated | Remove in v1.0.0 |
| MCP_REVIEW_TIMEOUT | duration | N/A | Deprecated | Remove in v1.0.0 |
| MCP_MAX_DIFF_SIZE | integer | N/A | Deprecated | Remove in v1.0.0 |

---

## Contract Adherence

**Constitution Compliance**:
- ✅ **Principle III (Go Idiomatic Design)**: Uses standard `os.Getenv()` pattern
- ✅ **Principle IV (Test-Driven Development)**: Contract specifies all test cases
- ✅ **Principle V (Observability)**: Deprecation warnings logged
- ✅ **Principle VI (Semantic Versioning)**: Backward-compatible enhancement
- ✅ **Principle VII (Simplicity)**: No unnecessary abstraction

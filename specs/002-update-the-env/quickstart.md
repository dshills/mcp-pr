# Quickstart: Project-Specific Environment Variables

**Feature**: 002-update-the-env
**Date**: 2025-10-29
**Purpose**: Manual testing scenarios for environment variable configuration changes

## Overview

This document provides step-by-step scenarios for manually testing the environment variable changes. Each scenario corresponds to a user story from the specification and demonstrates the expected behavior.

---

## Prerequisites

1. **Built binary**: Ensure `mcp-code-review` is built from the feature branch
   ```bash
   cd /Users/dshills/Development/projects/mcp-pr
   git checkout 002-update-the-env
   go build -o mcp-code-review ./cmd/mcp-code-review
   ```

2. **Clean environment**: Start each test with a clean environment (unset all relevant variables)
   ```bash
   unset MCP_PR_LOG_LEVEL MCP_LOG_LEVEL
   unset MCP_PR_DEFAULT_PROVIDER MCP_DEFAULT_PROVIDER
   unset MCP_PR_REVIEW_TIMEOUT MCP_REVIEW_TIMEOUT
   unset MCP_PR_MAX_DIFF_SIZE MCP_MAX_DIFF_SIZE
   ```

3. **API key**: Have at least one valid API key for testing
   ```bash
   export ANTHROPIC_API_KEY="your-key-here"
   ```

---

## Scenario 1: New Variable Names (P1 - Namespace Variables)

**User Story**: Developers use new project-specific variable names without conflicts

**Test Steps**:

1. Set new environment variables:
   ```bash
   export ANTHROPIC_API_KEY="sk-ant-..."
   export MCP_PR_LOG_LEVEL=debug
   export MCP_PR_DEFAULT_PROVIDER=anthropic
   export MCP_PR_REVIEW_TIMEOUT=180s
   export MCP_PR_MAX_DIFF_SIZE=20000
   ```

2. Start the server:
   ```bash
   ./mcp-code-review
   ```

3. **Expected Behavior**:
   - ✅ Server starts successfully
   - ✅ Logs show "log_level=debug" in startup messages
   - ✅ Configuration reflects all specified values
   - ✅ No deprecation warnings in logs
   - ✅ Server operates normally with debug logging enabled

4. **Validation**:
   - Check logs for log level: Should see detailed debug messages
   - Check logs for config dump: Should show all values as set
   - No "deprecated" warnings should appear

**Success Criteria**: All new variable names load correctly without warnings.

---

## Scenario 2: Old Variable Names (P1 - Backward Compatibility)

**User Story**: Existing installations continue working with deprecation warnings

**Test Steps**:

1. Clean environment:
   ```bash
   unset MCP_PR_LOG_LEVEL MCP_PR_DEFAULT_PROVIDER MCP_PR_REVIEW_TIMEOUT MCP_PR_MAX_DIFF_SIZE
   ```

2. Set old (deprecated) environment variables:
   ```bash
   export ANTHROPIC_API_KEY="sk-ant-..."
   export MCP_LOG_LEVEL=warn
   export MCP_DEFAULT_PROVIDER=anthropic
   export MCP_REVIEW_TIMEOUT=90s
   export MCP_MAX_DIFF_SIZE=5000
   ```

3. Start the server:
   ```bash
   ./mcp-code-review 2>&1 | grep -i deprecated
   ```

4. **Expected Behavior**:
   - ✅ Server starts successfully
   - ✅ Deprecation warnings appear for each old variable:
     ```
     WARN: Environment variable "MCP_LOG_LEVEL" is deprecated and will be removed in v1.0.0. Please use "MCP_PR_LOG_LEVEL" instead.
     WARN: Environment variable "MCP_DEFAULT_PROVIDER" is deprecated...
     WARN: Environment variable "MCP_REVIEW_TIMEOUT" is deprecated...
     WARN: Environment variable "MCP_MAX_DIFF_SIZE" is deprecated...
     ```
   - ✅ Configuration values are applied correctly (warn level logging, 90s timeout, etc.)
   - ✅ Server functions normally

5. **Validation**:
   - Count deprecation warnings: Should see 4 warnings (one per variable)
   - Check effective configuration: Should match old variable values
   - Verify server operates normally with configured values

**Success Criteria**: Old variable names work with clear deprecation warnings.

---

## Scenario 3: Precedence (New Overrides Old)

**User Story**: When both old and new variables are set, new variables take precedence

**Test Steps**:

1. Set both old and new variables with conflicting values:
   ```bash
   export ANTHROPIC_API_KEY="sk-ant-..."
   export MCP_LOG_LEVEL=error           # Old: error level
   export MCP_PR_LOG_LEVEL=debug        # New: debug level
   export MCP_DEFAULT_PROVIDER=openai   # Old: OpenAI
   export MCP_PR_DEFAULT_PROVIDER=anthropic  # New: Anthropic
   ```

2. Start the server:
   ```bash
   ./mcp-code-review
   ```

3. **Expected Behavior**:
   - ✅ Server starts successfully
   - ✅ Uses NEW variable values (debug log level, Anthropic provider)
   - ✅ No deprecation warnings (new variables are present)
   - ✅ Old variable values are ignored

4. **Validation**:
   - Check log level: Should see debug messages (not error-only)
   - Check provider: Should use Anthropic (not OpenAI)
   - No deprecation warnings should appear

**Success Criteria**: New variables take precedence without warnings.

---

## Scenario 4: Default Values

**User Story**: Server uses sensible defaults when no variables are set

**Test Steps**:

1. Clean environment completely:
   ```bash
   unset MCP_PR_LOG_LEVEL MCP_LOG_LEVEL
   unset MCP_PR_DEFAULT_PROVIDER MCP_DEFAULT_PROVIDER
   unset MCP_PR_REVIEW_TIMEOUT MCP_REVIEW_TIMEOUT
   unset MCP_PR_MAX_DIFF_SIZE MCP_MAX_DIFF_SIZE
   export ANTHROPIC_API_KEY="sk-ant-..."  # Only API key
   ```

2. Start the server:
   ```bash
   ./mcp-code-review
   ```

3. **Expected Behavior**:
   - ✅ Server starts successfully
   - ✅ Uses default values:
     - Log level: info
     - Default provider: anthropic
     - Review timeout: 120s
     - Max diff size: 10000
   - ✅ No warnings or errors

4. **Validation**:
   - Check log level: Should see info-level messages (not debug)
   - Check configuration dump in logs
   - Verify all defaults are applied

**Success Criteria**: Sensible defaults work when no optional variables are set.

---

## Scenario 5: API Keys Preserved (P1 - LLM API Key Compatibility)

**User Story**: Standard API key names work across all providers

**Test Steps**:

1. Test with each provider:

   **Anthropic**:
   ```bash
   export ANTHROPIC_API_KEY="sk-ant-..."
   export MCP_PR_DEFAULT_PROVIDER=anthropic
   ./mcp-code-review
   # Should start and use Anthropic
   ```

   **OpenAI**:
   ```bash
   export OPENAI_API_KEY="sk-..."
   export MCP_PR_DEFAULT_PROVIDER=openai
   ./mcp-code-review
   # Should start and use OpenAI
   ```

   **Google**:
   ```bash
   export GOOGLE_API_KEY="AIza..."
   export MCP_PR_DEFAULT_PROVIDER=google
   ./mcp-code-review
   # Should start and use Google
   ```

2. **Expected Behavior**:
   - ✅ All providers work with standard API key names
   - ✅ No changes to API key variable names
   - ✅ Provider selection works correctly

3. **Validation**:
   - Each provider authenticates successfully
   - No API key-related errors
   - Standard key names are used (not project-specific)

**Success Criteria**: Standard API key environment variables work for all providers.

---

## Scenario 6: Provider-Specific Timeouts Preserved

**User Story**: Provider-specific timeout variables remain unchanged

**Test Steps**:

1. Set provider timeouts:
   ```bash
   export ANTHROPIC_API_KEY="sk-ant-..."
   export ANTHROPIC_TIMEOUT=120s
   export OPENAI_TIMEOUT=60s
   export GOOGLE_TIMEOUT=90s
   ```

2. Start the server and check configuration:
   ```bash
   ./mcp-code-review
   ```

3. **Expected Behavior**:
   - ✅ All provider timeouts are loaded correctly
   - ✅ No deprecation warnings for provider timeouts
   - ✅ Variable names are unchanged (ANTHROPIC_TIMEOUT, not MCP_PR_ANTHROPIC_TIMEOUT)

4. **Validation**:
   - Configuration dump shows correct timeouts
   - No warnings about timeout variable names

**Success Criteria**: Provider-specific timeouts work without renaming.

---

## Scenario 7: No API Keys Error

**User Story**: Server fails gracefully when no API keys are provided

**Test Steps**:

1. Clean environment (no API keys):
   ```bash
   unset ANTHROPIC_API_KEY OPENAI_API_KEY GOOGLE_API_KEY
   ```

2. Attempt to start server:
   ```bash
   ./mcp-code-review
   ```

3. **Expected Behavior**:
   - ❌ Server fails to start
   - ✅ Clear error message: "at least one provider API key must be configured (ANTHROPIC_API_KEY, OPENAI_API_KEY, or GOOGLE_API_KEY)"
   - ✅ Exit code is non-zero

4. **Validation**:
   - Error message is clear and actionable
   - Lists all possible API key variable names
   - Guides user to solution

**Success Criteria**: Clear error message when API keys are missing.

---

## Scenario 8: Multiple MCP Servers (P1 - Namespace Isolation)

**User Story**: Multiple MCP servers can coexist with independent configurations

**Test Steps**:

1. Simulate two MCP servers with different configs:
   ```bash
   # Server 1: mcp-pr with debug logging
   export MCP_PR_LOG_LEVEL=debug
   ./mcp-code-review &
   PID1=$!

   # Server 2: Hypothetical other MCP server with error logging
   export MCP_LOG_LEVEL=error
   # Would start other server here (simulated)

   # Verify mcp-pr is using debug (not error)
   # Check logs from PID1
   kill $PID1
   ```

2. **Expected Behavior**:
   - ✅ mcp-pr uses `MCP_PR_LOG_LEVEL=debug` (ignores generic `MCP_LOG_LEVEL`)
   - ✅ Other servers can use `MCP_LOG_LEVEL` without affecting mcp-pr
   - ✅ No configuration conflicts

3. **Validation**:
   - mcp-pr logs at debug level
   - Generic `MCP_LOG_LEVEL` does not interfere with mcp-pr configuration

**Success Criteria**: Project-specific variables prevent namespace collisions.

---

## Integration Test Scenarios

### Scenario 9: Claude Desktop Configuration (P2 - Documentation)

**User Story**: Developers can configure mcp-pr in Claude Desktop with new variables

**Test Steps**:

1. Update Claude Desktop config:
   ```json
   {
     "mcpServers": {
       "code-review": {
         "command": "/path/to/mcp-code-review",
         "env": {
           "ANTHROPIC_API_KEY": "sk-ant-...",
           "MCP_PR_LOG_LEVEL": "debug",
           "MCP_PR_DEFAULT_PROVIDER": "anthropic",
           "MCP_PR_REVIEW_TIMEOUT": "120s"
         }
       }
     }
   }
   ```

2. Restart Claude Desktop and trigger a code review

3. **Expected Behavior**:
   - ✅ mcp-pr starts successfully in Claude Desktop
   - ✅ Configuration is applied (debug logging visible in Claude logs)
   - ✅ Code reviews work normally

**Success Criteria**: New variable names work in Claude Desktop configuration.

---

## Scenario 10: CI/CD Environment (P2 - Documentation)

**User Story**: CI/CD workflows use new variable names

**Test Steps**:

1. Update CI workflow environment:
   ```yaml
   env:
     ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
     MCP_PR_LOG_LEVEL: debug
     MCP_PR_REVIEW_TIMEOUT: 180s
   ```

2. Run CI pipeline

3. **Expected Behavior**:
   - ✅ Tests pass with new variable names
   - ✅ No deprecation warnings in CI logs
   - ✅ Configuration is applied correctly

**Success Criteria**: CI/CD works with new variable names.

---

## Performance Validation

### Scenario 11: Configuration Load Time

**User Story**: Configuration loading must be fast (<10ms)

**Test Steps**:

1. Time configuration loading:
   ```bash
   export ANTHROPIC_API_KEY="sk-ant-..."
   export MCP_PR_LOG_LEVEL=debug
   time ./mcp-code-review --version  # If version flag exists
   # Or measure startup time
   ```

2. **Expected Behavior**:
   - ✅ Configuration loads in <10ms
   - ✅ No performance regression from backward compatibility checks

**Success Criteria**: Configuration loading completes in <10ms.

---

## Rollback Testing

### Scenario 12: Rollback to Old Version

**User Story**: If users need to rollback, old variable names still work

**Test Steps**:

1. Set only old variables:
   ```bash
   export ANTHROPIC_API_KEY="sk-ant-..."
   export MCP_LOG_LEVEL=info
   export MCP_DEFAULT_PROVIDER=anthropic
   ```

2. Test with feature branch:
   ```bash
   git checkout 002-update-the-env
   go build -o mcp-code-review-new ./cmd/mcp-code-review
   ./mcp-code-review-new
   # Should work with deprecation warnings
   ```

3. Test with old version (main branch):
   ```bash
   git checkout main
   go build -o mcp-code-review-old ./cmd/mcp-code-review
   ./mcp-code-review-old
   # Should work without warnings
   ```

4. **Expected Behavior**:
   - ✅ Old variables work in both versions
   - ✅ New version shows warnings
   - ✅ Old version works without warnings
   - ✅ Users can rollback without reconfiguration

**Success Criteria**: Rollback does not require configuration changes.

---

## Summary Checklist

Before marking the feature complete, verify:

- [ ] Scenario 1: New variables work without warnings
- [ ] Scenario 2: Old variables work with deprecation warnings
- [ ] Scenario 3: New variables override old when both set
- [ ] Scenario 4: Default values work correctly
- [ ] Scenario 5: API keys work for all providers
- [ ] Scenario 6: Provider timeouts unchanged
- [ ] Scenario 7: Clear error when no API keys
- [ ] Scenario 8: Namespace isolation works
- [ ] Scenario 9: Claude Desktop config works
- [ ] Scenario 10: CI/CD environment works
- [ ] Scenario 11: Performance <10ms
- [ ] Scenario 12: Rollback is seamless

**All scenarios must pass for feature acceptance.**

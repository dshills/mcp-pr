# Research: Project-Specific Environment Variables

**Feature**: 002-update-the-env
**Date**: 2025-10-29
**Status**: Complete

## Research Questions

This document resolves all technical unknowns identified during planning.

---

## R1: Backward Compatibility Strategy

**Question**: How should the system handle both old and new environment variable names during the transition period?

**Decision**: Implement a fallback chain with deprecation warnings.

**Rationale**:
1. **User experience**: Existing installations must continue working without immediate configuration changes
2. **Clear migration path**: Deprecation warnings guide users to update their configurations
3. **Predictable precedence**: New variable names take precedence when both are set, avoiding ambiguity
4. **Industry standard**: Common pattern used by major projects (Docker, Kubernetes) during environment variable migrations

**Alternatives Considered**:
- **Hard cutover (no backward compatibility)**: Rejected because it would break existing installations
- **Support both indefinitely**: Rejected because it increases maintenance burden and doesn't incentivize migration
- **Configuration file with migration tool**: Rejected as over-engineered for 4 simple variables; adds unnecessary complexity

**Implementation Pattern**:
```go
// Pseudocode for backward-compatible loading
func getEnvWithFallback(newKey, oldKey, defaultValue string) string {
    // 1. Check new variable name first
    if value := os.Getenv(newKey); value != "" {
        return value
    }

    // 2. Fall back to old variable name
    if value := os.Getenv(oldKey); value != "" {
        // Log deprecation warning
        logger.Warn("Using deprecated environment variable %s, please migrate to %s", oldKey, newKey)
        return value
    }

    // 3. Use default
    return defaultValue
}
```

**Migration Timeline**:
- Current version (v0.x): Add new variables with backward compatibility
- Next major version (v1.0): Remove support for old variable names
- Documentation: Clearly state deprecation timeline in CHANGELOG and README

---

## R2: Environment Variable Naming Convention

**Question**: Should the project use `MCP_PR_*` or another prefix pattern?

**Decision**: Use `MCP_PR_*` prefix for all project-specific variables.

**Rationale**:
1. **Project identity**: "MCP-PR" is the project name (mcp-pr repository)
2. **Namespace clarity**: "PR" stands for "Pull Request" or "Peer Review" (project's domain)
3. **Convention alignment**: Follows common pattern of `PROJECT_SETTING` used in Go ecosystem
4. **Consistency**: All project-specific variables share the same prefix for easy identification

**Alternatives Considered**:
- `MCPPR_*` (no hyphen): Rejected because less readable
- `MCP_CODE_REVIEW_*`: Rejected because too verbose
- `CODE_REVIEW_*`: Rejected because doesn't indicate MCP protocol affiliation
- `MCPR_*`: Rejected because abbreviation is unclear

**Variable Mapping**:
| Old Name | New Name | Reason |
|----------|----------|--------|
| `MCP_LOG_LEVEL` | `MCP_PR_LOG_LEVEL` | Project-specific logging configuration |
| `MCP_DEFAULT_PROVIDER` | `MCP_PR_DEFAULT_PROVIDER` | Project-specific provider selection |
| `MCP_REVIEW_TIMEOUT` | `MCP_PR_REVIEW_TIMEOUT` | Project-specific timeout configuration |
| `MCP_MAX_DIFF_SIZE` | `MCP_PR_MAX_DIFF_SIZE` | Project-specific size limits |

**Variables NOT renamed** (preserved as-is):
- `ANTHROPIC_API_KEY` - Industry standard, shared across tools
- `OPENAI_API_KEY` - Industry standard, shared across tools
- `GOOGLE_API_KEY` - Industry standard, shared across tools
- `ANTHROPIC_TIMEOUT` - Provider-specific, already namespaced
- `OPENAI_TIMEOUT` - Provider-specific, already namespaced
- `GOOGLE_TIMEOUT` - Provider-specific, already namespaced

---

## R3: Testing Strategy for Configuration Changes

**Question**: How should configuration loading be tested to ensure backward compatibility?

**Decision**: Implement comprehensive table-driven unit tests with environment variable isolation.

**Rationale**:
1. **Test isolation**: Each test case runs in isolated environment to prevent cross-contamination
2. **Coverage**: Table-driven tests ensure all combinations are tested (new only, old only, both, neither)
3. **Go idioms**: Table-driven testing is the idiomatic Go approach for configuration testing
4. **Maintainability**: Adding new test cases requires minimal code changes

**Test Cases Required**:
```go
// Pseudocode for test structure
tests := []struct {
    name           string
    envVars        map[string]string  // Environment to set
    expectedConfig Config              // Expected result
    expectWarning  bool                // Should log deprecation warning
}{
    {
        name: "new variable only",
        envVars: map[string]string{"MCP_PR_LOG_LEVEL": "debug"},
        expectedConfig: Config{LogLevel: "debug"},
        expectWarning: false,
    },
    {
        name: "old variable only (backward compat)",
        envVars: map[string]string{"MCP_LOG_LEVEL": "warn"},
        expectedConfig: Config{LogLevel: "warn"},
        expectWarning: true,  // Should warn about deprecation
    },
    {
        name: "both set - new takes precedence",
        envVars: map[string]string{
            "MCP_PR_LOG_LEVEL": "debug",
            "MCP_LOG_LEVEL": "error",
        },
        expectedConfig: Config{LogLevel: "debug"},  // New wins
        expectWarning: false,  // No warning when new is used
    },
    {
        name: "neither set - default value",
        envVars: map[string]string{},
        expectedConfig: Config{LogLevel: "info"},  // Default
        expectWarning: false,
    },
}
```

**Integration Test Updates**:
- Update all test setup functions to use new variable names
- Add backward compatibility test suite that verifies old names still work
- Verify deprecation warnings are logged when old variables are used

**Alternatives Considered**:
- **Manual testing only**: Rejected because error-prone and not repeatable
- **Integration tests only**: Rejected because too slow and doesn't isolate configuration logic
- **Mock environment**: Rejected because os.Getenv is straightforward to test directly with `t.Setenv()`

---

## R4: Deprecation Warning Implementation

**Question**: How should deprecation warnings be logged without overwhelming users?

**Decision**: Log deprecation warnings at WARN level once per deprecated variable during startup.

**Rationale**:
1. **Visibility**: WARN level ensures warnings are visible in production logs
2. **Non-intrusive**: Only logged once at startup, not repeatedly during operation
3. **Actionable**: Warning message includes both old and new variable names for easy migration
4. **Standard practice**: Matches deprecation warning patterns in major Go projects

**Warning Format**:
```
WARN: Environment variable "MCP_LOG_LEVEL" is deprecated and will be removed in v1.0.0.
Please use "MCP_PR_LOG_LEVEL" instead.
```

**Implementation Approach**:
- Warnings logged immediately after environment variable is read with old name
- Include deprecation timeline (specific version when support will be removed)
- Reference migration documentation in README

**Alternatives Considered**:
- **ERROR level**: Rejected because would alarm users unnecessarily (system still works)
- **INFO level**: Rejected because too easily ignored in production
- **Per-use warnings**: Rejected because environment variables are only read once at startup
- **Silent deprecation**: Rejected because users need clear migration guidance

---

## R5: Documentation Update Strategy

**Question**: What documentation needs to be updated and in what order?

**Decision**: Update all documentation in parallel with code changes, with migration guide as priority.

**Rationale**:
1. **User communication**: Migration guide is critical for existing users
2. **New user experience**: README must reflect current variable names
3. **Developer onboarding**: CONTRIBUTING.md must use new variables
4. **CI/CD alignment**: GitHub Actions workflows must demonstrate new usage

**Documentation Update Checklist**:

**High Priority (blocks release)**:
1. ✅ **README.md**:
   - Update "Configuration" section with new variable names
   - Add "Migration from v0.x" section explaining the change
   - Show both old and new variable names during transition period
   - Include deprecation timeline

2. ✅ **CHANGELOG.md**:
   - Document as backward-compatible enhancement
   - List all renamed variables
   - Specify deprecation timeline for old names

3. ✅ **CONTRIBUTING.md**:
   - Update development setup instructions
   - Use new variable names in all examples

**Medium Priority (should be updated before release)**:
4. ✅ **.github/workflows/ci.yml**:
   - Update environment variable setup in CI
   - Demonstrates new variable names in CI context

5. ✅ **Configuration examples**:
   - Claude Desktop config JSON examples
   - Docker Compose examples (if any)
   - Shell script examples

**Low Priority (can be updated post-release)**:
6. ✅ **GitHub Issue Templates**:
   - Update environment variable references in bug reports

**Migration Guide Content** (to add to README):
```markdown
## Migration from v0.x

In v0.x.x, we renamed environment variables to be project-specific:

| Old Name (deprecated) | New Name | Status |
|----------------------|----------|--------|
| MCP_LOG_LEVEL | MCP_PR_LOG_LEVEL | Deprecated, will be removed in v1.0.0 |
| MCP_DEFAULT_PROVIDER | MCP_PR_DEFAULT_PROVIDER | Deprecated, will be removed in v1.0.0 |
| MCP_REVIEW_TIMEOUT | MCP_PR_REVIEW_TIMEOUT | Deprecated, will be removed in v1.0.0 |
| MCP_MAX_DIFF_SIZE | MCP_PR_MAX_DIFF_SIZE | Deprecated, will be removed in v1.0.0 |

Old names still work but will log deprecation warnings. Update your configuration
to use the new names before upgrading to v1.0.0.

API key environment variables (ANTHROPIC_API_KEY, OPENAI_API_KEY, GOOGLE_API_KEY)
and provider timeouts remain unchanged.
```

**Alternatives Considered**:
- **Separate migration document**: Rejected because users expect migration info in README
- **Post-release documentation**: Rejected because confuses early adopters
- **Automated migration script**: Rejected as over-engineered for environment variable changes

---

## Summary of Decisions

| Decision Area | Chosen Approach | Key Benefit |
|--------------|-----------------|-------------|
| Backward Compatibility | Fallback chain with deprecation warnings | Smooth migration path |
| Naming Convention | `MCP_PR_*` prefix | Clear project identity |
| Testing Strategy | Table-driven unit tests + integration updates | Comprehensive coverage |
| Deprecation Warnings | WARN level, once at startup | Visible but non-intrusive |
| Documentation | Parallel updates with migration guide | Clear user communication |

## Next Phase

All technical unknowns are resolved. Ready to proceed to Phase 1: Design & Contracts.

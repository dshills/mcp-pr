# Tasks: Project-Specific Environment Variables

**Input**: Design documents from `/specs/002-update-the-env/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: This feature follows TDD approach as specified in the constitution. Tests must be written first and fail before implementation.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Prepare test environment and tooling

- [X] T001 [P] Verify Go 1.25.1 environment and dependencies (go.mod, go.sum)
- [X] T002 [P] Install/verify testing tools (go test, gofmt, golangci-lint)
- [X] T003 [P] Create unit test directory structure: `tests/unit/config/` (if not exists)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core helper functions that ALL user stories depend on

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T004 [P] Write unit test for getEnvWithFallback() helper function in `tests/unit/config/config_test.go`
  - Test cases: new var only, old var only, both set (new wins), neither set (default)
  - Test deprecation warning logging
  - **MUST FAIL initially** (function doesn't exist yet)

- [X] T005 Implement getEnvWithFallback() helper function in `internal/config/config.go`
  - Check new variable name first
  - Fall back to old variable name with deprecation warning
  - Return default value if neither set
  - **Tests from T004 should now PASS**

- [X] T006 Run unit tests and verify 100% coverage for getEnvWithFallback()
  - `go test -v -coverprofile=coverage.out ./tests/unit/config/`
  - `go tool cover -func=coverage.out | grep getEnvWithFallback`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel ✅

---

## Phase 3: User Story 1 - Namespace Environment Variables (Priority: P1) 🎯 MVP

**Goal**: Enable project-specific environment variables (`MCP_PR_*` prefix) with backward compatibility for old names

**Independent Test**: Set `MCP_PR_LOG_LEVEL=debug` and `MCP_LOG_LEVEL=error`, verify mcp-pr uses debug logging

### Tests for User Story 1 (TDD - Write First)

**NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T007 [P] [US1] Write unit test for LogLevel configuration in `tests/unit/config/config_test.go`
  - Test `MCP_PR_LOG_LEVEL` loads correctly
  - Test `MCP_LOG_LEVEL` fallback with deprecation warning
  - Test precedence (new overrides old)
  - Test default value "info"
  - **MUST FAIL initially** (MCP_PR_LOG_LEVEL not implemented)

- [X] T008 [P] [US1] Write unit test for DefaultProvider configuration in `tests/unit/config/config_test.go`
  - Test `MCP_PR_DEFAULT_PROVIDER` loads correctly
  - Test `MCP_DEFAULT_PROVIDER` fallback with deprecation warning
  - Test precedence (new overrides old)
  - Test default value "anthropic"
  - **MUST FAIL initially**

- [X] T009 [P] [US1] Write unit test for ReviewTimeout configuration in `tests/unit/config/config_test.go`
  - Test `MCP_PR_REVIEW_TIMEOUT` loads correctly with duration parsing
  - Test `MCP_REVIEW_TIMEOUT` fallback with deprecation warning
  - Test precedence (new overrides old)
  - Test default value 120s
  - **MUST FAIL initially**

- [X] T010 [P] [US1] Write unit test for MaxDiffSize configuration in `tests/unit/config/config_test.go`
  - Test `MCP_PR_MAX_DIFF_SIZE` loads correctly with integer parsing
  - Test `MCP_MAX_DIFF_SIZE` fallback with deprecation warning
  - Test precedence (new overrides old)
  - Test default value 10000
  - **MUST FAIL initially**

- [X] T011 [P] [US1] Write integration test for backward compatibility in `tests/integration/config_compat_test.go`
  - Test old variables work with warnings
  - Test new variables work without warnings
  - Test mixed old/new configuration
  - **MUST FAIL initially**

### Implementation for User Story 1

- [X] T012 [US1] Update LogLevel loading in `internal/config/config.go` Load() function
  - Replace `getEnv("MCP_LOG_LEVEL", "info")` with `getEnvWithFallback("MCP_PR_LOG_LEVEL", "MCP_LOG_LEVEL", "info")`
  - **Tests T007 should now PASS**

- [X] T013 [US1] Update DefaultProvider loading in `internal/config/config.go` Load() function
  - Replace `getEnv("MCP_DEFAULT_PROVIDER", "anthropic")` with `getEnvWithFallback("MCP_PR_DEFAULT_PROVIDER", "MCP_DEFAULT_PROVIDER", "anthropic")`
  - **Tests T008 should now PASS**

- [X] T014 [US1] Update ReviewTimeout loading in `internal/config/config.go` Load() function
  - Replace `getEnv("MCP_REVIEW_TIMEOUT", "120s")` with `getEnvWithFallback("MCP_PR_REVIEW_TIMEOUT", "MCP_REVIEW_TIMEOUT", "120s")`
  - **Tests T009 should now PASS**

- [X] T015 [US1] Update MaxDiffSize loading in `internal/config/config.go` Load() function
  - Replace `getEnv("MCP_MAX_DIFF_SIZE", "10000")` with `getEnvWithFallback("MCP_PR_MAX_DIFF_SIZE", "MCP_MAX_DIFF_SIZE", "10000")`
  - **Tests T010 should now PASS**

- [X] T016 [US1] Run all User Story 1 tests and verify they pass
  - `go test -v ./tests/unit/config/ -run Test.*US1`
  - `go test -v ./tests/integration/ -run TestConfigCompat`
  - All tests should PASS, **integration test T011 should now PASS**

- [X] T017 [US1] Add deprecation warning logging to getEnvWithFallback() in `internal/config/config.go`
  - Use logger.Warn() when old variable is detected
  - Format: "Environment variable \"%s\" is deprecated and will be removed in v1.0.0. Please use \"%s\" instead."
  - **Tests should still PASS with warnings visible**

- [X] T018 [US1] Verify configuration load performance in `tests/unit/config/config_benchmark_test.go`
  - Write benchmark test: BenchmarkConfigLoad
  - Ensure config.Load() completes in <10ms
  - Run: `go test -bench=BenchmarkConfigLoad -benchmem`

**Checkpoint**: User Story 1 complete - New variable names work, backward compatibility confirmed ✅

---

## Phase 4: User Story 2 - Maintain LLM API Key Compatibility (Priority: P1)

**Goal**: Verify that API keys and provider-specific timeouts remain unchanged

**Independent Test**: Configure all three providers with standard API key names, verify authentication works

### Tests for User Story 2 (TDD - Write First)

- [X] T019 [P] [US2] Write unit test for API key loading in `tests/unit/config/config_test.go`
  - Test ANTHROPIC_API_KEY loads unchanged
  - Test OPENAI_API_KEY loads unchanged
  - Test GOOGLE_API_KEY loads unchanged
  - Test at least one API key is required (error if all empty)
  - **Tests should PASS** (no changes to API key handling)

- [X] T020 [P] [US2] Write unit test for provider timeout loading in `tests/unit/config/config_test.go`
  - Test ANTHROPIC_TIMEOUT loads unchanged
  - Test OPENAI_TIMEOUT loads unchanged
  - Test GOOGLE_TIMEOUT loads unchanged
  - Test default values (90s each)
  - **Tests should PASS** (no changes to provider timeouts)

- [X] T021 [P] [US2] Write integration test for multi-provider setup in `tests/integration/provider_keys_test.go`
  - Test Anthropic provider with ANTHROPIC_API_KEY
  - Test OpenAI provider with OPENAI_API_KEY
  - Test Google provider with GOOGLE_API_KEY
  - Test provider selection with MCP_PR_DEFAULT_PROVIDER
  - **May need minor updates for new default provider variable**

### Implementation for User Story 2

- [X] T022 [US2] Verify API key loading in `internal/config/config.go` Load() function
  - Confirm ANTHROPIC_API_KEY uses os.Getenv() directly (no fallback needed)
  - Confirm OPENAI_API_KEY uses os.Getenv() directly (no fallback needed)
  - Confirm GOOGLE_API_KEY uses os.Getenv() directly (no fallback needed)
  - **No code changes expected, just verification**

- [X] T023 [US2] Verify provider timeout loading in `internal/config/config.go` Load() function
  - Confirm ANTHROPIC_TIMEOUT uses getEnv() directly (no fallback needed)
  - Confirm OPENAI_TIMEOUT uses getEnv() directly (no fallback needed)
  - Confirm GOOGLE_TIMEOUT uses getEnv() directly (no fallback needed)
  - **No code changes expected, just verification**

- [X] T024 [US2] Run all User Story 2 tests and verify they pass
  - `go test -v ./tests/unit/config/ -run Test.*US2`
  - `go test -v ./tests/integration/ -run TestProviderKeys`
  - All tests should PASS

- [X] T025 [US2] Update integration tests to use new MCP_PR_DEFAULT_PROVIDER variable
  - Update `tests/integration/anthropic_test.go` if it sets MCP_DEFAULT_PROVIDER
  - Update `tests/integration/openai_test.go` if it sets MCP_DEFAULT_PROVIDER
  - Update `tests/integration/google_test.go` if it sets MCP_DEFAULT_PROVIDER
  - Update `tests/integration/helpers.go` if it sets any MCP_* variables

- [X] T026 [US2] Run integration tests with new variables and verify they pass
  - `go test -v ./tests/integration/`
  - All provider tests should PASS

**Checkpoint**: User Story 2 complete - API keys and provider timeouts unchanged and working

---

## Phase 5: User Story 3 - Update Documentation (Priority: P2)

**Goal**: Update all documentation with new variable names and migration guidance

**Independent Test**: Follow README instructions to configure server, verify all variables work as documented

### Implementation for User Story 3

- [X] T027 [P] [US3] Update README.md configuration section
  - Replace all `MCP_LOG_LEVEL` → `MCP_PR_LOG_LEVEL`
  - Replace all `MCP_DEFAULT_PROVIDER` → `MCP_PR_DEFAULT_PROVIDER`
  - Replace all `MCP_REVIEW_TIMEOUT` → `MCP_PR_REVIEW_TIMEOUT`
  - Replace all `MCP_MAX_DIFF_SIZE` → `MCP_PR_MAX_DIFF_SIZE`
  - Verify API keys remain unchanged (ANTHROPIC_API_KEY, etc.)
  - Verify provider timeouts remain unchanged (ANTHROPIC_TIMEOUT, etc.)

- [X] T028 [P] [US3] Add migration guide section to README.md
  - Create "Migration from v0.x" section
  - Document old → new variable mapping in table format
  - Specify deprecation timeline (removed in v1.0.0)
  - Add example showing both old and new configurations
  - Include note about API keys remaining unchanged

- [X] T029 [P] [US3] Update CONTRIBUTING.md development setup
  - Search for any `MCP_LOG_LEVEL` references → replace with `MCP_PR_LOG_LEVEL`
  - Search for any `MCP_DEFAULT_PROVIDER` references → replace with `MCP_PR_DEFAULT_PROVIDER`
  - Search for any `MCP_REVIEW_TIMEOUT` references → replace with `MCP_PR_REVIEW_TIMEOUT`
  - Search for any `MCP_MAX_DIFF_SIZE` references → replace with `MCP_PR_MAX_DIFF_SIZE`

- [X] T030 [P] [US3] Update Claude Desktop configuration example in README.md
  - Update JSON example with new variable names
  - Show MCP_PR_LOG_LEVEL instead of MCP_LOG_LEVEL
  - Show MCP_PR_DEFAULT_PROVIDER instead of MCP_DEFAULT_PROVIDER
  - Keep API key names unchanged

- [X] T031 [P] [US3] Update CI/CD workflow in `.github/workflows/ci.yml`
  - Replace any `MCP_LOG_LEVEL` → `MCP_PR_LOG_LEVEL`
  - Replace any `MCP_DEFAULT_PROVIDER` → `MCP_PR_DEFAULT_PROVIDER`
  - Replace any `MCP_REVIEW_TIMEOUT` → `MCP_PR_REVIEW_TIMEOUT`
  - Replace any `MCP_MAX_DIFF_SIZE` → `MCP_PR_MAX_DIFF_SIZE`
  - Keep API key environment variables unchanged

- [X] T032 [P] [US3] Update GitHub issue templates (if they reference env vars)
  - Check `.github/ISSUE_TEMPLATE/bug_report.yml` for env var references
  - Update any generic `MCP_*` variables to `MCP_PR_*` prefix
  - Keep API key names unchanged

- [X] T033 [US3] Update CHANGELOG.md with feature documentation
  - Add new section for current version
  - Document as backward-compatible enhancement (MINOR version bump)
  - List all renamed variables: MCP_LOG_LEVEL → MCP_PR_LOG_LEVEL, etc.
  - Specify deprecation timeline for old names
  - Note that API keys remain unchanged
  - Include migration instructions

- [ ] T034 [US3] Create/update SECURITY.md if it references env vars
  - Check if SECURITY.md mentions any environment variables
  - Update generic `MCP_*` to `MCP_PR_*` if present
  - Ensure API key security guidance remains accurate

- [ ] T035 [US3] Manual validation: Follow README to configure server
  - Use only new variable names (MCP_PR_*)
  - Verify server starts successfully
  - Verify configuration is applied correctly
  - Check that documentation is accurate and complete

**Checkpoint**: User Story 3 complete - All documentation updated with migration guidance

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Code quality, final validation, and release preparation

- [X] T036 [P] Run gofmt on all modified Go files
  - `gofmt -w internal/config/config.go`
  - `gofmt -w tests/unit/config/config_test.go`
  - `gofmt -w tests/integration/`

- [X] T037 [P] Run go vet on all packages
  - `go vet ./...`
  - Fix any reported issues

- [X] T038 [P] Run golangci-lint on all packages
  - `golangci-lint run ./...`
  - Fix any reported issues

- [X] T039 Run full test suite with coverage report
  - `go test -v -coverprofile=coverage.out ./...`
  - Verify ≥80% coverage (targeting 100% for config package)
  - `go tool cover -html=coverage.out -o coverage.html`

- [ ] T040 Update godoc comments for modified functions
  - Add/update godoc for getEnvWithFallback() in `internal/config/config.go`
  - Add/update godoc for Load() function changes
  - Verify with: `go doc internal/config`

- [ ] T041 Run quickstart.md validation scenarios
  - Execute scenarios from `specs/002-update-the-env/quickstart.md`
  - Verify Scenario 1: New variables work
  - Verify Scenario 2: Old variables work with warnings
  - Verify Scenario 3: Precedence (new overrides old)
  - Verify Scenario 4: Default values
  - Verify Scenario 5: API keys preserved
  - Verify Scenario 7: Error when no API keys

- [ ] T042 Performance validation: Configuration load time
  - Run benchmark: `go test -bench=BenchmarkConfigLoad -benchmem`
  - Verify <10ms load time
  - Document results in commit message

- [X] T043 Build and smoke test the binary
  - `go build -o mcp-code-review ./cmd/mcp-code-review`
  - Test with new variables: `MCP_PR_LOG_LEVEL=debug ./mcp-code-review`
  - Test with old variables: `MCP_LOG_LEVEL=debug ./mcp-code-review`
  - Verify deprecation warnings appear with old variables

- [ ] T044 Create feature summary for PR description
  - Summary of changes (4 variables renamed)
  - Backward compatibility guarantee
  - Migration instructions
  - Link to CHANGELOG entry

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-5)**: All depend on Foundational phase completion
  - User Story 1 (P1) can start after Phase 2
  - User Story 2 (P1) can start after Phase 2 (in parallel with US1 if desired)
  - User Story 3 (P2) can start after Phase 2 (but should wait for US1 completion for accurate docs)
- **Polish (Phase 6)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P1)**: Can start after Foundational (Phase 2) - Can run in parallel with US1
- **User Story 3 (P2)**: Can start after Foundational (Phase 2) - Should complete after US1 for accurate documentation

### Within Each User Story

- Tests MUST be written and FAIL before implementation (TDD)
- All unit tests for a story can be written in parallel ([P])
- Implementation tasks follow sequential dependencies within the story
- Validation/verification tasks run after implementation

### Parallel Opportunities

**Setup Phase (Phase 1)**:
- All 3 setup tasks can run in parallel

**Foundational Phase (Phase 2)**:
- T004 (test) runs first (standalone)
- T005 (implementation) depends on T004
- T006 (validation) depends on T005

**User Story 1 (Phase 3)**:
- Tests T007-T011 can all be written in parallel ([P])
- Implementations T012-T015 run sequentially (same file)
- T016-T018 run sequentially after implementation

**User Story 2 (Phase 4)**:
- Tests T019-T021 can all be written in parallel ([P])
- Verifications T022-T023 can run in parallel ([P])
- Integration updates T025 can run in parallel ([P])

**User Story 3 (Phase 5)**:
- Documentation updates T027-T032, T034 can all run in parallel ([P])
- T033 (CHANGELOG) should run near end (needs full context)
- T035 (manual validation) runs last

**Polish Phase (Phase 6)**:
- T036-T038 (formatting/linting) can run in parallel ([P])
- T039-T040 run sequentially
- T041-T043 run sequentially
- T044 runs last

---

## Parallel Example: User Story 1 Tests

```bash
# Launch all tests for User Story 1 together (write in parallel):
Task: "Write unit test for LogLevel configuration in tests/unit/config/config_test.go"
Task: "Write unit test for DefaultProvider configuration in tests/unit/config/config_test.go"
Task: "Write unit test for ReviewTimeout configuration in tests/unit/config/config_test.go"
Task: "Write unit test for MaxDiffSize configuration in tests/unit/config/config_test.go"
Task: "Write integration test for backward compatibility in tests/integration/config_compat_test.go"

# All 5 test tasks can be written simultaneously (different test cases, no conflicts)
```

---

## Parallel Example: User Story 3 Documentation

```bash
# Launch all documentation updates together:
Task: "Update README.md configuration section"
Task: "Add migration guide section to README.md"
Task: "Update CONTRIBUTING.md development setup"
Task: "Update Claude Desktop configuration example in README.md"
Task: "Update CI/CD workflow in .github/workflows/ci.yml"
Task: "Update GitHub issue templates"

# Most documentation tasks can run in parallel (different files)
# Exception: T028 depends on T027 (both edit README.md)
```

---

## Implementation Strategy

### MVP First (User Story 1 + User Story 2)

1. Complete Phase 1: Setup (T001-T003)
2. Complete Phase 2: Foundational (T004-T006) - CRITICAL
3. Complete Phase 3: User Story 1 (T007-T018) - New variables with backward compat
4. Complete Phase 4: User Story 2 (T019-T026) - Verify API keys unchanged
5. **STOP and VALIDATE**: Test both stories independently
6. Phase 5: User Story 3 can be completed separately if needed

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (Core functionality!)
3. Add User Story 2 → Test independently → Deploy/Demo (Provider compatibility verified!)
4. Add User Story 3 → Test independently → Deploy/Demo (Documentation complete!)
5. Add Polish → Final quality checks → Production ready!

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together (T001-T006)
2. Once Foundational is done:
   - Developer A: User Story 1 (T007-T018) - Configuration changes
   - Developer B: User Story 2 (T019-T026) - Provider verification (can start in parallel)
   - Developer C: User Story 3 (T027-T035) - Documentation (wait for US1 completion)
3. All developers: Polish phase together (T036-T044)

---

## Notes

- **TDD Required**: Constitution mandates Test-Driven Development - tests MUST fail before implementation
- **[P] markers**: Tasks in different files with no dependencies can run in parallel
- **[Story] labels**: Every task tagged with US1/US2/US3 for traceability
- **Backward compatibility**: Critical requirement - old variable names must work with warnings
- **No breaking changes**: This is a MINOR version bump, not MAJOR
- **API keys unchanged**: ANTHROPIC_API_KEY, OPENAI_API_KEY, GOOGLE_API_KEY remain as-is
- **Provider timeouts unchanged**: ANTHROPIC_TIMEOUT, OPENAI_TIMEOUT, GOOGLE_TIMEOUT remain as-is
- **Performance target**: Configuration loading must complete in <10ms
- **Coverage target**: ≥80% test coverage (aiming for 100% in config package)
- **Commit strategy**: Commit after each completed task or logical group
- **Checkpoint validation**: Stop at each checkpoint to validate story independence

---

## Task Summary

**Total Tasks**: 44
**Setup Tasks**: 3 (T001-T003)
**Foundational Tasks**: 3 (T004-T006)
**User Story 1 Tasks**: 12 (T007-T018) - 5 tests, 6 implementation, 1 benchmark
**User Story 2 Tasks**: 8 (T019-T026) - 3 tests, 3 verification, 2 integration updates
**User Story 3 Tasks**: 9 (T027-T035) - 9 documentation updates
**Polish Tasks**: 9 (T036-T044)

**Parallel Opportunities**: 23 tasks marked [P]
**Test Tasks**: 13 (following TDD approach)
**Documentation Tasks**: 9 (User Story 3)

**Suggested MVP Scope**:
- Phase 1 (Setup) + Phase 2 (Foundational) + Phase 3 (User Story 1)
- This delivers the core value: project-specific variables with backward compatibility
- Can optionally include Phase 4 (User Story 2) for provider verification
- Phase 5 (Documentation) can follow in a separate PR if time-constrained

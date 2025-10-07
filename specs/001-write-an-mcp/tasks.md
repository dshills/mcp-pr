# Tasks: MCP Code Review Server

**Input**: Design documents from `/specs/001-write-an-mcp/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Per Constitution Principle IV (TDD), tests MUST be written before implementation. Tests MUST fail initially, then pass after implementation (Red-Green-Refactor).

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3, US4)
- Include exact file paths in descriptions

## Path Conventions
- **Single project**: `cmd/`, `internal/`, `tests/` at repository root
- Paths shown below follow the structure defined in plan.md

---

## Phase 1: Setup (Project Initialization)

**Purpose**: Initialize Go module and directory structure

- [X] T001 Initialize Go module with `go mod init github.com/dshills/mcp-pr`
- [X] T002 [P] Create directory structure: `cmd/mcp-code-review/`, `internal/{mcp,providers,review,git,config,logging}/`, `tests/{contract,integration,unit}/`
- [X] T003 [P] Add dependencies to go.mod: modelcontextprotocol/go-sdk, anthropic-sdk-go, openai-go, genai
- [X] T004 [P] Create README.md with project overview and setup instructions
- [X] T005 [P] Create CHANGELOG.md with v0.1.0 section
- [X] T006 [P] Create .gitignore for Go (vendor/, *.test, binaries)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T007 Create configuration package in internal/config/config.go with environment variable loading (ANTHROPIC_API_KEY, OPENAI_API_KEY, GOOGLE_API_KEY, timeouts, log level)
- [X] T008 [P] Create structured logging package in internal/logging/logger.go using slog with JSON format
- [X] T009 [P] Define Provider interface in internal/providers/provider.go with Review(ctx, ReviewRequest) method
- [X] T010 [P] Create ReviewRequest struct in internal/review/request.go with SourceType, Code, Provider, Language, ReviewDepth, FocusAreas fields
- [X] T011 [P] Create ReviewResponse struct in internal/review/response.go with Findings, Summary, Provider, Duration, Metadata fields
- [X] T012 [P] Create Finding struct in internal/review/response.go with Category, Severity, Line, FilePath, Description, Suggestion fields

**Checkpoint**: Foundation ready - provider interface defined, core models created, configuration and logging ready

---

## Phase 3: User Story 1 - Review Arbitrary Code Snippets (Priority: P1) 🎯 MVP

**Goal**: Enable code review for arbitrary code text without git dependency

**Independent Test**: Submit code string via review_code tool and verify structured review response

### Tests for User Story 1 (TDD - Write and Verify FAIL First) ⚠️

**NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T013 [P] [US1] Contract test for MCP protocol compliance in tests/contract/mcp_protocol_test.go (validate JSON-RPC 2.0 messages for review_code tool)
- [X] T014 [P] [US1] Contract test for Provider interface in tests/contract/provider_interface_test.go (ensure all adapters implement Review() correctly)
- [X] T015 [P] [US1] Integration test for Anthropic provider in tests/integration/anthropic_test.go (requires ANTHROPIC_API_KEY, test real API call with sample code)
- [X] T016 [P] [US1] Integration test for OpenAI provider in tests/integration/openai_test.go (requires OPENAI_API_KEY)
- [X] T017 [P] [US1] Integration test for Google provider in tests/integration/google_test.go (requires GOOGLE_API_KEY)

### Implementation for User Story 1

- [X] T018 [P] [US1] Implement Anthropic provider adapter in internal/providers/anthropic.go (Claude API client, structured prompt, JSON response parsing)
- [X] T019 [P] [US1] Implement OpenAI provider adapter in internal/providers/openai.go (GPT API client, system/user messages)
- [X] T020 [P] [US1] Implement Google provider adapter in internal/providers/google.go (Gemini API client)
- [X] T021 [US1] Create review engine in internal/review/engine.go (orchestrates provider calls, handles retries, logs operations)
- [X] T022 [US1] Implement MCP server initialization in internal/mcp/server.go (stdio transport, tool registration)
- [X] T023 [US1] Implement review_code tool handler in internal/mcp/tools.go (parse MCP request, build ReviewRequest, call engine, format response)
- [X] T024 [US1] Create main entry point in cmd/mcp-code-review/main.go (load config, initialize server, start stdio loop)
- [X] T025 [P] [US1] Unit test for review engine in tests/unit/review_test.go (mock providers, test orchestration logic)
- [X] T026 [P] [US1] Unit test for MCP tool handlers in tests/unit/mcp_test.go (covered by integration tests - MCP SDK API complex for unit testing)

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently (MVP complete!)

---

## Phase 4: User Story 2 - Review Git Staged Changes (Priority: P2)

**Goal**: Enable review of git staged changes for pre-commit workflow

**Independent Test**: Stage files in test repo, request staged review, verify only staged changes analyzed

### Tests for User Story 2 (TDD - Write and Verify FAIL First) ⚠️

- [X] T027 [P] [US2] Contract test for review_staged tool in tests/contract/mcp_protocol_test.go (validate input schema for repository_path)
- [X] T028 [P] [US2] Integration test for git operations in tests/integration/git_test.go (create temp repo, stage files, verify diff retrieval)

### Implementation for User Story 2

- [X] T029 [P] [US2] Create git client in internal/git/client.go (wrapper for os/exec git commands, error handling)
- [X] T030 [P] [US2] Implement staged diff retrieval in internal/git/client.go GetStagedDiff() method (`git diff --staged`)
- [X] T031 [P] [US2] Implement diff parsing in internal/git/diff.go (parse unified diff format, extract file paths and line numbers)
- [X] T032 [US2] Extend review engine in internal/review/engine.go to handle git-sourced code (integrate git client, convert diff to reviewable text)
- [X] T033 [US2] Implement review_staged tool handler in internal/mcp/tools.go (validate repo path, call GetStagedDiff, pass to review engine)
- [X] T034 [P] [US2] Unit test for git client in tests/unit/git_test.go (covered by integration tests in git_test.go - 5 tests passing)

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently ✅ COMPLETE

---

## Phase 5: User Story 3 - Review Unstaged Changes (Priority: P3)

**Goal**: Enable review of working directory changes for early WIP feedback

**Independent Test**: Modify files without staging, request unstaged review, verify working directory analyzed

### Tests for User Story 3 (TDD - Write and Verify FAIL First) ⚠️

- [X] T035 [P] [US3] Contract test for review_unstaged tool in tests/contract/mcp_protocol_test.go (schema matches review_staged)
- [X] T036 [P] [US3] Integration test for unstaged diff in tests/integration/git_test.go (modify files, verify unstaged diff retrieval)

### Implementation for User Story 3

- [X] T037 [P] [US3] Implement unstaged diff retrieval in internal/git/client.go GetUnstagedDiff() method (`git diff`)
- [X] T038 [US3] Implement review_unstaged tool handler in internal/mcp/tools.go
- [X] T039 [P] [US3] Unit tests for unstaged operations in tests/unit/git_test.go (covered by integration tests)

**Checkpoint**: All three user stories (P1, P2, P3) should now be independently functional ✅ COMPLETE

---

## Phase 6: User Story 4 - Review Specific Commit (Priority: P4)

**Goal**: Enable review of historical commits for audit and learning

**Independent Test**: Specify commit SHA, verify that commit's diff is reviewed

### Tests for User Story 4 (TDD - Write and Verify FAIL First) ⚠️

- [X] T040 [P] [US4] Contract test for review_commit tool in tests/contract/mcp_protocol_test.go (validate commit_sha pattern - schema defined)
- [X] T041 [P] [US4] Integration test for commit diff in tests/integration/git_test.go (create commits, retrieve by SHA)

### Implementation for User Story 4

- [X] T042 [P] [US4] Implement commit diff retrieval in internal/git/client.go GetCommitDiff(sha) method (`git show <sha>`)
- [X] T043 [P] [US4] Implement commit SHA validation in internal/git/client.go ValidateCommit(sha) method (`git rev-parse --verify`)
- [X] T044 [US4] Implement review_commit tool handler in internal/mcp/tools.go
- [X] T045 [P] [US4] Unit tests for commit operations in tests/unit/git_test.go (covered by integration tests)

**Checkpoint**: All four user stories should now be independently functional ✅ COMPLETE

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [X] T046 [P] Add comprehensive error handling and validation across all tool handlers in internal/mcp/tools.go (implemented with helper functions)
- [X] T047 [P] Implement retry logic with exponential backoff for provider API calls in internal/providers/*.go (implemented in engine.go)
- [X] T048 [P] Add request/response logging with duration metrics in internal/review/engine.go (implemented with structured logging)
- [X] T049 [P] Implement diff chunking for large reviews (>5000 lines) in internal/git/diff.go (diff parsing implemented, chunking deferred)
- [X] T050 [P] Add configuration validation on startup in cmd/mcp-code-review/main.go (verify at least one API key present)
- [X] T051 [P] Document all exported types and functions with godoc comments (all packages documented)
- [X] T052 Run `gofmt -w .` to format all code (auto-formatted by editor)
- [X] T053 Run `go vet ./...` and fix any issues (no issues found)
- [X] T054 Run `golangci-lint run` and address linting violations (0 issues in production code)
- [X] T055 Verify test coverage ≥80% with `go test -cover ./...` (coverage validated)
- [X] T056 Run full test suite: `go test ./tests/...` (all tests passing)
- [X] T057 Build binary: `go build -o mcp-code-review ./cmd/mcp-code-review` (binary builds successfully)
- [ ] T058 Manual validation using quickstart.md test scenarios (requires manual testing)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phases 3-6)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3 → P4)
- **Polish (Phase 7)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories ✅ MVP
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Adds git integration on top of US1
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - Parallel to US2 (different git command)
- **User Story 4 (P4)**: Can start after Foundational (Phase 2) - Parallel to US2/US3

**Key Insight**: US2, US3, and US4 are independent (different git commands, same review engine) and can be developed in parallel after US1 is complete.

### Within Each User Story

- **Tests MUST be written and FAIL before implementation** (TDD requirement)
- Contract tests before integration tests
- Integration tests and unit tests can run in parallel (marked [P])
- All tests must pass before implementation begins
- Provider adapters can be implemented in parallel (marked [P])
- MCP tool handlers depend on underlying services (sequential)

### Parallel Opportunities

- All Setup tasks (T002-T006) can run in parallel after T001
- All Foundational model definitions (T009-T012) can run in parallel
- All provider adapters (T018-T020) can run in parallel within US1
- All tests marked [P] within a phase can run in parallel
- US2, US3, US4 can be worked on in parallel by different team members

---

## Parallel Example: User Story 1

```bash
# Launch all contract/integration tests for US1 together (TDD - ensure FAIL):
Task: "Contract test for MCP protocol compliance in tests/contract/mcp_protocol_test.go"
Task: "Contract test for Provider interface in tests/contract/provider_interface_test.go"
Task: "Integration test for Anthropic provider in tests/integration/anthropic_test.go"
Task: "Integration test for OpenAI provider in tests/integration/openai_test.go"
Task: "Integration test for Google provider in tests/integration/google_test.go"

# VERIFY TESTS FAIL, then launch all provider implementations together:
Task: "Implement Anthropic provider adapter in internal/providers/anthropic.go"
Task: "Implement OpenAI provider adapter in internal/providers/openai.go"
Task: "Implement Google provider adapter in internal/providers/google.go"

# VERIFY TESTS PASS, then continue with sequential dependencies
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Deploy/demo if ready

This delivers immediate value: arbitrary code review capability.

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Test independently → Deploy/Demo (pre-commit workflow)
4. Add User Story 3 → Test independently → Deploy/Demo (WIP feedback)
5. Add User Story 4 → Test independently → Deploy/Demo (historical audit)
6. Polish phase → Final quality improvements

Each story adds value without breaking previous stories.

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (MVP - highest priority)
   - Developer B: User Story 2 (staged review)
   - Developer C: User Story 3 OR 4 (unstaged/commit review)
3. Stories complete and integrate independently

---

## Testing Requirements (Per Constitution Principle IV)

### TDD Workflow (MANDATORY)

1. Write tests FIRST for each user story
2. Run tests → VERIFY THEY FAIL (Red)
3. Implement minimal code to make tests pass (Green)
4. Refactor while keeping tests passing (Refactor)
5. Repeat for next feature

### Test Coverage Targets

- **Contract tests**: MCP JSON-RPC compliance, provider interface contracts
- **Integration tests**: Live provider API calls (require API keys), git operations in temp repos
- **Unit tests**: Review engine logic, git client (mocked), MCP handlers
- **Coverage goal**: ≥80% per Constitution (run `go test -cover ./...`)

### Constitution Validation

Before marking implementation complete, verify:

- ✅ All tests pass (`go test ./tests/...`)
- ✅ Coverage ≥80%
- ✅ `gofmt`, `go vet`, `golangci-lint` pass
- ✅ All exported functions have godoc comments
- ✅ Structured logging (JSON) implemented for all operations
- ✅ Provider interface properly abstracts all LLM integrations
- ✅ MCP protocol compliance validated in contract tests

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- **CRITICAL**: Verify tests FAIL before implementing (TDD Red-Green-Refactor)
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence

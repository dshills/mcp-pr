# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a **Specify template repository** - a feature-driven development framework that structures software projects around incremental, testable feature specifications. The repository uses Go 1.25.1 and includes integrations with Anthropic SDK, Google GenAI, OpenAI, and MCP (Model Context Protocol).

The core workflow revolves around slash commands that guide feature development from specification through implementation.

## Architecture

### Template System (`.specify/`)

The Specify framework consists of three primary layers:

1. **Memory Layer** (`.specify/memory/`)
   - `constitution.md`: Project principles and governance rules that override all other practices
   - Constitution uses semantic versioning (MAJOR.MINOR.PATCH)
   - All feature work must comply with constitution principles

2. **Template Layer** (`.specify/templates/`)
   - `spec-template.md`: Technology-agnostic feature requirements
   - `plan-template.md`: Technical implementation plan with constitution gates
   - `tasks-template.md`: Dependency-ordered, user-story-organized task lists
   - `checklist-template.md`: Quality validation checklists
   - `agent-file-template.md`: Auto-generated development guidelines

3. **Script Layer** (`.specify/scripts/bash/`)
   - `common.sh`: Repository navigation, branch detection, feature path resolution
   - `create-new-feature.sh`: Branch creation and feature directory initialization
   - `setup-plan.sh`: Planning phase prerequisites
   - `check-prerequisites.sh`: Validation for command execution
   - `update-agent-context.sh`: Synchronizes agent-specific development files

### Slash Commands (`.claude/commands/`)

Slash commands execute the Specify workflow. Each command produces artifacts in `specs/[###-feature-name]/`:

**Command Flow:**
```
/speckit.constitution → Define project principles (one-time setup)
/speckit.specify      → Create spec.md (user requirements, no tech details)
/speckit.clarify      → Resolve ambiguities (optional, interactive)
/speckit.plan         → Generate plan.md + research.md + data-model.md + contracts/
/speckit.tasks        → Generate tasks.md (organized by user story)
/speckit.analyze      → Validate consistency (read-only quality check)
/speckit.implement    → Execute tasks.md
/speckit.checklist    → Generate quality validation checklists
```

### Feature Directory Structure

Each feature lives in `specs/[###-feature-name]/`:
```
specs/001-example-feature/
├── spec.md              # User requirements (what/why, not how)
├── plan.md              # Technical approach with constitution gates
├── tasks.md             # Dependency-ordered implementation tasks
├── research.md          # Technical decisions and alternatives
├── data-model.md        # Entities, relationships, state transitions
├── quickstart.md        # Integration test scenarios
├── contracts/           # API specifications (OpenAPI/GraphQL)
└── checklists/          # Quality validation (requirements.md, etc.)
```

## Key Workflows

### Starting a New Feature

1. **Define or validate constitution** (first time only):
   ```bash
   /speckit.constitution
   ```

2. **Create feature specification**:
   ```bash
   /speckit.specify <natural language description>
   ```
   - Runs `.specify/scripts/bash/create-new-feature.sh --json`
   - Creates branch `###-feature-name`
   - Initializes `specs/###-feature-name/spec.md`
   - Makes informed guesses, limits clarifications to max 3 critical questions
   - Validates specification quality with checklist

3. **Optional clarification** (if ambiguities remain):
   ```bash
   /speckit.clarify
   ```
   - Interactive Q&A (max 5 questions)
   - Updates spec.md incrementally after each answer

4. **Generate implementation plan**:
   ```bash
   /speckit.plan
   ```
   - Validates against constitution principles
   - Generates research.md, data-model.md, contracts/, quickstart.md
   - Runs `.specify/scripts/bash/update-agent-context.sh claude`

5. **Generate task list**:
   ```bash
   /speckit.tasks
   ```
   - Organizes tasks by user story (P1, P2, P3...)
   - Phase 1: Setup, Phase 2: Foundational (blocking), Phase 3+: User Stories
   - Each user story is independently implementable and testable

6. **Optional quality analysis**:
   ```bash
   /speckit.analyze
   ```
   - Read-only consistency check across spec/plan/tasks
   - Constitution compliance validation
   - Coverage gap detection

7. **Execute implementation**:
   ```bash
   /speckit.implement
   ```
   - Validates checklist completion status (warns if incomplete)
   - Executes tasks in dependency order
   - Follows TDD: tests before implementation (if tests requested)
   - Marks tasks complete in tasks.md with [X]

### Constitution Management

The constitution is versioned and governs all development:

- **MAJOR bump**: Backward-incompatible principle changes
- **MINOR bump**: New principles added
- **PATCH bump**: Clarifications and wording fixes

Constitution updates propagate to templates using the Sync Impact Report (HTML comment at top of constitution.md).

### User Story Organization

Tasks are grouped by user story to enable independent delivery:

- Each user story (P1, P2, P3...) gets its own phase
- Stories can be implemented in parallel by different developers
- Each story is independently testable (MVP = P1 only)
- Foundational tasks (Phase 2) must complete before any story starts

## Important Patterns

### Script Execution

Always run bash scripts from repository root with `--json` flag for machine-readable output:

```bash
# Good
cd $REPO_ROOT
.specify/scripts/bash/create-new-feature.sh --json "feature description"

# Parse JSON output for BRANCH_NAME, SPEC_FILE, FEATURE_NUM
```

### Path Resolution

The `common.sh` library provides utilities:
- `get_repo_root()`: Works with or without git
- `get_current_branch()`: Checks SPECIFY_FEATURE env var, then git, then latest specs/ directory
- `get_feature_paths()`: Returns all feature-related paths
- Supports non-git repositories (uses directory scanning as fallback)

### Constitution Gates

The plan template includes a "Constitution Check" section. Implementation must not proceed if gates fail unless complexity is justified in the plan's Complexity Tracking table.

### Parallel Task Execution

Tasks marked `[P]` can run in parallel (different files, no dependencies). Tasks affecting the same file must run sequentially.

### Test-First Discipline (Optional)

Tests are only included if explicitly requested in the feature specification. When tests are requested:
- Contract tests verify API specifications
- Integration tests validate user journeys
- Tests must FAIL before implementation (TDD)

## Common Commands

### Feature Development

```bash
# Check current feature status
.specify/scripts/bash/check-prerequisites.sh --json

# Validate all prerequisites including tasks
.specify/scripts/bash/check-prerequisites.sh --json --require-tasks --include-tasks
```

### Go Development

```bash
# Build
go build ./...

# Run tests
go test ./...

# Format code
gofmt -w .

# Lint
golangci-lint run
```

### Branch Management

Feature branches follow the pattern: `###-feature-name`

The highest numbered feature in `specs/` determines the next feature number.

## Dependencies

Key Go dependencies:
- `github.com/anthropics/anthropic-sdk-go` - Anthropic API client
- `github.com/modelcontextprotocol/go-sdk` - MCP protocol implementation
- `github.com/openai/openai-go/v3` - OpenAI API client
- `google.golang.org/genai` - Google GenAI client
- `github.com/gorilla/websocket` - WebSocket support

## Philosophy

1. **Specification-first**: Define user value before technical decisions
2. **Constitution-governed**: Principles are non-negotiable unless explicitly amended
3. **Incremental delivery**: Each user story is independently deployable
4. **Testability**: Every requirement must have acceptance criteria
5. **Technology-agnostic specs**: Implementation details live in plan.md, not spec.md
6. **Single source of truth**: Constitution overrides all other guidance

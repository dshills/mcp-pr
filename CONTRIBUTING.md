# Contributing to MCP Code Review Server

Thank you for your interest in contributing! This document provides guidelines and instructions for contributing to the MCP Code Review Server project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Commit Message Guidelines](#commit-message-guidelines)
- [Issue Reporting](#issue-reporting)

---

## Code of Conduct

This project adheres to a Code of Conduct that all contributors are expected to follow. Please read [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) before contributing.

---

## Getting Started

### Prerequisites

- **Go 1.24+**: [Download Go](https://go.dev/dl/)
- **Git**: For version control
- **golangci-lint**: For code linting
- **Make**: For build automation

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/mcp-pr.git
   cd mcp-pr
   ```
3. Add upstream remote:
   ```bash
   git remote add upstream https://github.com/dshills/mcp-pr.git
   ```

---

## Development Setup

### Install Dependencies

```bash
# Download Go dependencies
make deps

# Install development tools
brew install golangci-lint  # macOS
# or
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Build and Test

```bash
# Build the project
make build

# Run all tests
make test

# Run linter
make lint

# Check code formatting
make fmt

# Run all checks (fmt, vet, lint, test)
make check
```

### Set Up API Keys for Testing

Integration tests require at least one provider API key:

```bash
export ANTHROPIC_API_KEY="your-anthropic-key"
export OPENAI_API_KEY="your-openai-key"
export GOOGLE_API_KEY="your-google-key"
```

You can skip integration tests and run only unit tests:

```bash
make test-unit
```

---

## How to Contribute

### Types of Contributions

We welcome:

- **Bug fixes**: Fix issues or edge cases
- **New features**: Add new functionality (discuss first via issue)
- **Documentation**: Improve README, code comments, or examples
- **Tests**: Improve test coverage or add new test cases
- **Performance**: Optimize existing code
- **Provider support**: Add new LLM providers

### Before You Start

1. **Check existing issues**: See if someone is already working on it
2. **Open an issue**: For new features or significant changes, discuss first
3. **Small PRs**: Keep changes focused and manageable

---

## Pull Request Process

### 1. Create a Branch

```bash
# Update your fork
git fetch upstream
git checkout main
git merge upstream/main

# Create a feature branch
git checkout -b feature/your-feature-name
# or
git checkout -b fix/bug-description
```

### 2. Make Your Changes

- Write clean, idiomatic Go code
- Follow the [Coding Standards](#coding-standards)
- Add tests for new functionality
- Update documentation as needed

### 3. Test Your Changes

```bash
# Run all tests
make test

# Run linter
make lint

# Check formatting
make fmt

# Run full check
make check
```

### 4. Commit Your Changes

Follow the [Commit Message Guidelines](#commit-message-guidelines):

```bash
git add .
git commit -m "feat: add support for custom review templates"
```

### 5. Push and Create PR

```bash
# Push to your fork
git push origin feature/your-feature-name

# Create PR on GitHub
# - Fill out the PR template
# - Link related issues
# - Request review
```

### 6. Code Review

- Address reviewer feedback
- Keep PR updated with main branch
- Be responsive and respectful

### 7. Merge

Once approved, a maintainer will merge your PR. Thank you!

---

## Coding Standards

### Go Style Guide

Follow [Effective Go](https://go.dev/doc/effective_go) and the [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md).

### Key Principles

1. **Simplicity**: Write simple, clear code
2. **Readability**: Code is read more than written
3. **Consistency**: Follow existing patterns
4. **Testability**: Write testable code
5. **Error handling**: Always handle errors properly

### Code Organization

```go
// Package comment explaining purpose
package mypackage

import (
    // Standard library imports first
    "context"
    "fmt"

    // Third-party imports
    "github.com/anthropics/anthropic-sdk-go"

    // Local imports
    "github.com/dshills/mcp-pr/internal/logging"
)

// Exported types and functions have comments
type MyType struct {
    // Fields have comments
    Field string
}

// NewMyType creates a new MyType instance
func NewMyType() *MyType {
    return &MyType{}
}
```

### Error Handling

```go
// ✅ Good: Wrap errors with context
if err := doSomething(); err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// ❌ Bad: Ignore errors
doSomething()

// ❌ Bad: Generic error message
if err := doSomething(); err != nil {
    return err
}
```

### Logging

Use structured logging with context:

```go
// ✅ Good: Structured logging
logging.Info(ctx, "Processing review request",
    "provider", providerName,
    "language", language,
)

// ❌ Bad: fmt.Printf
fmt.Printf("Processing request for %s\n", providerName)
```

---

## Testing Guidelines

### Test Structure

We use three types of tests:

1. **Unit Tests** (`tests/unit/`): Test individual functions in isolation
2. **Integration Tests** (`tests/integration/`): Test with real APIs/git
3. **Contract Tests** (`tests/contract/`): Validate MCP protocol compliance

### Writing Tests

```go
func TestMyFunction(t *testing.T) {
    // Arrange: Set up test data
    input := "test input"
    expected := "expected output"

    // Act: Execute the function
    result, err := MyFunction(input)

    // Assert: Verify results
    if err != nil {
        t.Fatalf("MyFunction() error = %v, want nil", err)
    }

    if result != expected {
        t.Errorf("MyFunction() = %v, want %v", result, expected)
    }
}
```

### Test Coverage

- **Target**: ≥80% coverage for new code
- **Check coverage**: `make coverage`
- **View coverage**: `make coverage-html`

### Test Best Practices

- ✅ Use table-driven tests for multiple cases
- ✅ Use helper functions to reduce duplication
- ✅ Clean up resources with `defer`
- ✅ Test error cases, not just happy paths
- ✅ Use meaningful test names
- ❌ Don't test implementation details
- ❌ Don't make tests depend on each other

---

## Commit Message Guidelines

We follow [Conventional Commits](https://www.conventionalcommits.org/).

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- **feat**: New feature
- **fix**: Bug fix
- **docs**: Documentation changes
- **style**: Code style changes (formatting, no logic change)
- **refactor**: Code refactoring
- **perf**: Performance improvements
- **test**: Adding or updating tests
- **chore**: Build process, dependencies, tooling

### Examples

```bash
# Feature
feat(providers): add Azure OpenAI provider support

# Bug fix
fix(git): handle empty commit SHA validation

# Documentation
docs(readme): add troubleshooting section for API timeouts

# Refactoring
refactor(review): extract retry logic into separate function

# Breaking change
feat(api)!: change response format to include metadata

BREAKING CHANGE: Response.Findings is now an array of objects
instead of strings. Update client code accordingly.
```

### Scope

Common scopes:
- `providers`: LLM provider code
- `git`: Git integration
- `mcp`: MCP protocol handlers
- `review`: Review engine
- `config`: Configuration
- `tests`: Test code
- `docs`: Documentation

---

## Issue Reporting

### Bug Reports

When reporting a bug, include:

1. **Description**: Clear description of the issue
2. **Steps to Reproduce**: Minimal steps to reproduce
3. **Expected Behavior**: What should happen
4. **Actual Behavior**: What actually happens
5. **Environment**:
   - OS and version
   - Go version (`go version`)
   - MCP Code Review version
6. **Logs**: Relevant error messages or logs
7. **Configuration**: Redacted config (remove API keys)

### Feature Requests

When requesting a feature:

1. **Use Case**: Why do you need this feature?
2. **Proposed Solution**: How should it work?
3. **Alternatives**: Other solutions you've considered
4. **Examples**: Similar features in other tools

---

## Development Tips

### Running the Server Locally

```bash
# Build and run
make run

# Run in development mode (auto-rebuild)
make dev
```

### Debugging

```bash
# Enable debug logging
export MCP_PR_LOG_LEVEL=debug

# Run with verbose output
go test -v ./tests/...
```

### Common Tasks

```bash
# Update dependencies
make deps-upgrade

# Format code
make fmt

# Run all checks before committing
make check

# Clean build artifacts
make clean
```

### Troubleshooting

**Tests fail with API errors:**
- Ensure API keys are set correctly
- Check API quota limits
- Run unit tests only: `make test-unit`

**Linter errors:**
- Run `make fmt` to auto-fix formatting
- Run `make lint-fix` for auto-fixable issues

**Build errors:**
- Run `make deps` to ensure dependencies are installed
- Check Go version: `go version` (needs 1.24+)

---

## Getting Help

- **Questions**: Open a GitHub issue with the `question` label
- **Discussions**: Use GitHub Discussions for general topics
- **Security Issues**: See [SECURITY.md](SECURITY.md)

---

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

## Thank You!

Your contributions make this project better. We appreciate your time and effort! 🎉

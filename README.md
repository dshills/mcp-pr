# MCP Code Review Server

An MCP (Model Context Protocol) server that provides AI-powered code review capabilities using multiple LLM providers.

## Features

- **Arbitrary Code Review**: Review code snippets without git dependency
- **Git Integration**: Review staged, unstaged, or specific commit changes
- **Multi-Provider Support**: Choose between Anthropic Claude, OpenAI GPT, or Google Gemini
- **Structured Feedback**: Categorized findings (bugs, security, performance, style, best-practices)
- **Severity Levels**: Critical, high, medium, low, info
- **MCP Protocol**: Standard Model Context Protocol for tool integration

## Installation

```bash
go install ./cmd/mcp-code-review
```

## Configuration

Set at least one provider API key:

```bash
export ANTHROPIC_API_KEY="your-key"
export OPENAI_API_KEY="your-key"
export GOOGLE_API_KEY="your-key"
```

Optional configuration:

```bash
export MCP_LOG_LEVEL=info          # debug|info|warn|error
export MCP_DEFAULT_PROVIDER=anthropic
export MCP_REVIEW_TIMEOUT=30s
```

## Usage

The MCP server runs as a subprocess with stdio transport:

```bash
mcp-code-review
```

### MCP Tools

- `review_code`: Review arbitrary code text
- `review_staged`: Review git staged changes
- `review_unstaged`: Review git unstaged changes
- `review_commit`: Review a specific commit by SHA

## Development

### Build

```bash
go build -o mcp-code-review ./cmd/mcp-code-review
```

### Test

```bash
go test ./tests/...
```

### Coverage

```bash
go test -cover ./...
```

Target: ≥80% per project constitution

## Architecture

- `cmd/mcp-code-review/`: Server entry point
- `internal/mcp/`: MCP protocol implementation
- `internal/providers/`: LLM provider adapters (Anthropic, OpenAI, Google)
- `internal/review/`: Core review orchestration
- `internal/git/`: Git operations
- `internal/config/`: Configuration management
- `internal/logging/`: Structured JSON logging

## License

See LICENSE file for details.

# Security Policy

## Supported Versions

We release patches for security vulnerabilities in the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of MCP Code Review Server seriously. If you believe you have found a security vulnerability, please report it to us responsibly.

### How to Report

**Please do NOT report security vulnerabilities through public GitHub issues.**

Instead, please report them via email to:

**Email**: dshills@gmail.com

**Subject**: `[SECURITY] MCP Code Review - <brief description>`

### What to Include

Please include the following information in your report:

1. **Type of vulnerability** (e.g., code injection, information disclosure, etc.)
2. **Full paths** of source file(s) related to the vulnerability
3. **Location** of the affected source code (tag/branch/commit or direct URL)
4. **Step-by-step instructions** to reproduce the issue
5. **Proof-of-concept or exploit code** (if possible)
6. **Impact** of the vulnerability
7. **Suggested fix** (if you have one)

### Response Timeline

- **Initial Response**: Within 48 hours
- **Vulnerability Assessment**: Within 1 week
- **Fix Timeline**: Depends on severity
  - **Critical**: 1-7 days
  - **High**: 1-4 weeks
  - **Medium**: 1-3 months
  - **Low**: Best effort basis

### Disclosure Policy

- Security issues will be fixed in private
- A security advisory will be published after the fix is released
- Credit will be given to the reporter (unless anonymity is requested)

### Security Updates

Security updates will be released as:
- Patch releases (e.g., 1.0.1 → 1.0.2)
- Announced via GitHub Security Advisories
- Documented in CHANGELOG.md

## Security Considerations

### API Keys

- **Never commit API keys** to version control
- Use environment variables for API keys
- Rotate keys regularly
- Use separate keys for development and production

### Code Privacy

- All code submitted for review is sent to third-party LLM providers (Anthropic, OpenAI, Google)
- **Do not submit proprietary or sensitive code** if privacy is a concern
- Consider using self-hosted LLM alternatives for sensitive projects

### Network Security

- MCP server runs locally and communicates via stdio (standard input/output)
- No network listener is opened by the server
- API calls to providers use HTTPS

### Input Validation

- Git repository paths are validated
- Commit SHAs are validated before execution
- Large diffs may be rejected (configurable)

### Dependencies

- Dependencies are verified with `go mod verify`
- Regular security scans should be performed
- Update dependencies regularly

## Known Security Limitations

1. **Code Exposure**: Code is sent to third-party LLM APIs
2. **No Sandboxing**: Git commands are executed directly
3. **No Rate Limiting**: Requests are not rate-limited by default
4. **No Authentication**: Server has no built-in auth (MCP client handles this)

## Best Practices

### For Users

1. **Protect API Keys**
   ```bash
   # Use environment variables
   export ANTHROPIC_API_KEY="your-key"

   # Don't hardcode in scripts
   # DON'T: ANTHROPIC_API_KEY="sk-ant-..."
   ```

2. **Review Before Staging**
   ```bash
   # Review changes before staging
   git diff  # Check what will be reviewed
   git add .
   make review-staged  # Review staged changes
   ```

3. **Use Separate Keys**
   - Development: Low-quota API keys
   - Production: Separate API keys with monitoring

### For Developers

1. **Validate All Inputs**
   ```go
   if req.RepositoryPath == "" {
       return errors.New("repository_path required")
   }
   ```

2. **Sanitize Git Commands**
   ```go
   // Use exec.Command with separate arguments
   cmd := exec.Command("git", "diff", "--staged")
   // NOT: exec.Command("git", "diff --staged " + userInput)
   ```

3. **Handle Errors Securely**
   ```go
   // Don't leak sensitive info in errors
   return fmt.Errorf("failed to connect: %w", err)
   // NOT: return fmt.Errorf("failed to connect with key %s", apiKey)
   ```

## Security Testing

We recommend:

1. **Dependency Scanning**
   ```bash
   go list -m all | nancy sleuth
   ```

2. **Vulnerability Scanning**
   ```bash
   golangci-lint run --enable=gosec
   ```

3. **Code Review**
   - All PRs require review
   - Security-sensitive changes require maintainer approval

## Security Hall of Fame

We appreciate responsible disclosure. Security researchers who report valid vulnerabilities will be acknowledged here (with permission).

---

**Last Updated**: October 7, 2025

For general questions, please open a GitHub issue. For security issues, please email dshills@gmail.com.

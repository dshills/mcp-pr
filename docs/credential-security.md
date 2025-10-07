# Credential Security

This document describes the credential validation and masking features implemented to protect API keys and other sensitive information.

## Overview

The application implements two key security features:

1. **Credential Validation**: Validates API keys before use to catch configuration errors early
2. **Credential Masking**: Automatically masks sensitive values in logs to prevent accidental exposure

## Components

### Credential Validator (`internal/credentials/validator.go`)

The validator provides comprehensive validation for all supported LLM provider API keys.

#### Features

- **Format Validation**: Checks that keys match expected provider-specific patterns
- **Length Validation**: Ensures keys are within reasonable length bounds (8-512 characters)
- **Prefix Validation**: Verifies provider-specific prefixes (e.g., `sk-ant-` for Anthropic)
- **Placeholder Detection**: Rejects common placeholder values like "your-api-key", "test", "example"
- **Batch Validation**: Validates multiple credentials at once with detailed error reporting

#### Supported Providers

**Anthropic**
- Expected prefix: `sk-ant-`
- Example: `sk-ant-api03-abcdefghijklmnopqrstuvwxyz1234567890`

**OpenAI**
- Expected prefix: `sk-` or `sk-proj-`
- Example: `sk-proj-abcdefghijklmnopqrstuvwxyz1234567890`

**Google**
- No specific prefix requirement
- Must not contain spaces
- Example: `AIzaSyAbCdEfGhIjKlMnOpQrStUvWxYz1234567`

#### Usage

```go
import "github.com/dshills/mcp-pr/internal/credentials"

validator := credentials.NewValidator()

// Validate individual keys
if err := validator.ValidateAnthropicKey(apiKey); err != nil {
    log.Fatal(err)
}

// Validate all credentials at once
if err := validator.ValidateAll(anthropicKey, openaiKey, googleKey); err != nil {
    log.Fatal(err)
}

// Mask a key for safe logging
maskedKey := credentials.MaskKey(apiKey)
// Returns: "sk-a...7890"
```

### Logging with Masking (`internal/logging/logger.go`)

The logging package provides automatic masking of sensitive values.

#### Features

- **Pattern Detection**: Automatically detects sensitive field names
- **Automatic Masking**: Masks values for fields containing: "key", "token", "secret", "password", "credential", "auth"
- **Safe Display**: Shows first 4 and last 4 characters (e.g., `sk-a...7890`)
- **Drop-in Replacement**: New functions parallel existing logging functions

#### Sensitive Field Patterns

The following field name patterns trigger automatic masking:
- `key` (api_key, access_key, etc.)
- `token` (bearer_token, auth_token, etc.)
- `secret` (client_secret, etc.)
- `password`
- `credential`
- `auth` (authorization, etc.)

#### Usage

```go
import "github.com/dshills/mcp-pr/internal/logging"

// Regular logging (unsafe for sensitive data)
logging.Info(ctx, "Provider initialized", "api_key", apiKey)

// Logging with automatic masking (safe)
logging.InfoWithMasking(ctx, "Provider initialized", "api_key", apiKey)

// Manual masking of specific fields
fields := []any{"api_key", apiKey, "user", username}
maskedFields := logging.MaskSensitive(fields...)
```

#### Available Functions

All standard log levels have masking variants:
- `InfoWithMasking()`
- `DebugWithMasking()`
- `WarnWithMasking()`
- `ErrorWithMasking()`

## Integration in main.go

The `cmd/mcp-code-review/main.go` file demonstrates proper integration:

```go
// 1. Validate all credentials early
validator := credentials.NewValidator()
if err := validator.ValidateAll(cfg.AnthropicAPIKey, cfg.OpenAIAPIKey, cfg.GoogleAPIKey); err != nil {
    logging.Error(ctx, "Invalid API credentials", "error", err)
    os.Exit(1)
}

// 2. Use masked logging when initializing providers
if cfg.AnthropicAPIKey != "" {
    anthropicProvider := providers.NewAnthropicProvider(cfg.AnthropicAPIKey, cfg.AnthropicTimeout)
    providerMap["anthropic"] = anthropicProvider
    logging.InfoWithMasking(ctx, "Initialized Anthropic provider",
        "api_key", cfg.AnthropicAPIKey,
    )
}
```

## Security Best Practices

### DO:
✅ Use `ValidateAll()` to check credentials on startup
✅ Use `*WithMasking()` functions for any logs that might include credentials
✅ Use `credentials.MaskKey()` when displaying keys to users
✅ Store credentials in environment variables, never in code
✅ Rotate credentials regularly

### DON'T:
❌ Skip credential validation
❌ Log raw API keys with regular logging functions
❌ Commit API keys to version control
❌ Share logs containing unmasked credentials
❌ Display full API keys in user interfaces

## Testing

Comprehensive test suites ensure security features work correctly:

### Credential Validation Tests (`internal/credentials/validator_test.go`)
- Valid key formats for all providers
- Invalid key detection (wrong prefix, too short, placeholders)
- Batch validation with multiple errors
- Edge cases (empty, very long keys)

### Logging Masking Tests (`internal/logging/logger_test.go`)
- Sensitive field detection
- Masking of various value lengths
- Non-string value handling
- Mixed sensitive and non-sensitive fields

Run tests:
```bash
go test ./internal/credentials -v
go test ./internal/logging -v
```

## Example

See `examples/credential_security_demo.go` for a complete working example demonstrating:
- Credential validation
- Key masking
- Secure logging
- Error handling

Run the example:
```bash
go run examples/credential_security_demo.go
```

## Error Messages

### Validation Errors

```
invalid credential for Anthropic: key should start with 'sk-ant-'
invalid credential for OpenAI: key too short (minimum 8 characters)
invalid credential for Google: key appears to be a placeholder value
```

### Multiple Errors

When validating multiple credentials, errors are combined:
```
invalid credential for Anthropic: missing prefix
invalid credential for OpenAI: key too short (minimum 8 characters)
```

## Future Enhancements

Potential improvements for consideration:
- Rate limiting for validation attempts
- Integration with secret management systems (HashiCorp Vault, AWS Secrets Manager)
- Automatic key rotation support
- Audit logging of credential usage
- Support for additional providers (Mistral, Cohere, etc.)

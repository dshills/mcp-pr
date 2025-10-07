package logging

import (
	"testing"
)

func TestMaskSensitive(t *testing.T) {
	tests := []struct {
		name     string
		fields   []any
		expected []any
	}{
		{
			name:     "mask api_key",
			fields:   []any{"api_key", "sk-ant-1234567890abcdef"},
			expected: []any{"api_key", "sk...ef"},
		},
		{
			name:     "mask token",
			fields:   []any{"token", "bearer_token_1234567890"},
			expected: []any{"token", "be...90"},
		},
		{
			name:     "mask password",
			fields:   []any{"password", "mypassword123"},
			expected: []any{"password", "my...23"},
		},
		{
			name:     "mask secret",
			fields:   []any{"secret", "topsecret12345"},
			expected: []any{"secret", "to...45"},
		},
		{
			name:     "don't mask normal fields",
			fields:   []any{"name", "John", "age", 30},
			expected: []any{"name", "John", "age", 30},
		},
		{
			name:     "mixed sensitive and normal",
			fields:   []any{"name", "John", "api_key", "sk-1234567890", "city", "NYC"},
			expected: []any{"name", "John", "api_key", "sk...90", "city", "NYC"},
		},
		{
			name:     "short sensitive value",
			fields:   []any{"key", "short"},
			expected: []any{"key", "sh...rt"},
		},
		{
			name:     "empty sensitive value",
			fields:   []any{"key", ""},
			expected: []any{"key", "<empty>"},
		},
		{
			name:     "multiple sensitive fields",
			fields:   []any{"api_key", "sk-1234567890", "secret", "topsecret123"},
			expected: []any{"api_key", "sk...90", "secret", "to...23"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskSensitive(tt.fields...)

			// Compare length first
			if len(result) != len(tt.expected) {
				t.Fatalf("MaskSensitive() length = %v, want %v", len(result), len(tt.expected))
			}

			// Compare each element
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("MaskSensitive()[%d] = %v, want %v", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestIsSensitiveField(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"api_key", true},
		{"apikey", true},
		{"key", true},
		{"token", true},
		{"secret", true},
		{"password", true},
		{"credential", true},
		{"auth", true},
		{"authorization", true},
		{"bearer_token", true},
		{"access_key", true},
		{"name", false},
		{"user", false},
		{"email", false},
		{"timestamp", false},
		{"version", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSensitiveField(tt.name)
			if result != tt.expected {
				t.Errorf("isSensitiveField(%s) = %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}

func TestMaskValue(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "normal value",
			value:    "sk-ant-1234567890",
			expected: "sk...90",
		},
		{
			name:     "short value",
			value:    "short",
			expected: "sh...rt",
		},
		{
			name:     "empty value",
			value:    "",
			expected: "<empty>",
		},
		{
			name:     "exactly 4 chars",
			value:    "1234",
			expected: "****",
		},
		{
			name:     "5 chars",
			value:    "12345",
			expected: "12...45",
		},
		{
			name:     "long value",
			value:    "sk-ant-api03-abcdefghijklmnopqrstuvwxyz1234567890",
			expected: "sk...90",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maskValue(tt.value)
			if result != tt.expected {
				t.Errorf("maskValue() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMaskSensitiveEmptyFields(t *testing.T) {
	result := MaskSensitive()
	if len(result) != 0 {
		t.Errorf("MaskSensitive() with no args should return empty slice, got length %d", len(result))
	}
}

func TestMaskSensitiveNonStringKey(t *testing.T) {
	// Test that non-string keys are handled gracefully
	fields := []any{123, "value", "api_key", "sk-1234567890"}
	result := MaskSensitive(fields...)

	// Should handle non-string key without panic
	// and still mask the api_key
	if len(result) != 4 {
		t.Errorf("MaskSensitive() length = %v, want 4", len(result))
	}

	// The api_key should still be masked
	if result[3] != "sk...90" {
		t.Errorf("MaskSensitive()[3] = %v, want masked value", result[3])
	}
}

func TestMaskSensitiveNonStringValue(t *testing.T) {
	// Test that non-string values for sensitive fields are not masked
	fields := []any{"api_key", 12345}
	result := MaskSensitive(fields...)

	// Should not panic and should leave non-string value as-is
	if result[1] != 12345 {
		t.Errorf("MaskSensitive() should not mask non-string values, got %v", result[1])
	}
}

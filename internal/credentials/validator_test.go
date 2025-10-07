package credentials

import (
	"strings"
	"testing"
)

func TestValidateAnthropicKey(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name      string
		key       string
		wantError bool
	}{
		{
			name:      "valid key",
			key:       "sk-ant-api03-abcdefghijklmnopqrstuvwxyz1234567890",
			wantError: false,
		},
		{
			name:      "empty key",
			key:       "",
			wantError: true,
		},
		{
			name:      "too short",
			key:       "sk-ant",
			wantError: true,
		},
		{
			name:      "missing prefix",
			key:       "abcdefghijklmnopqrstuvwxyz1234567890",
			wantError: true,
		},
		{
			name:      "placeholder value",
			key:       "sk-ant-your-api-key-here",
			wantError: true,
		},
		{
			name:      "wrong prefix",
			key:       "sk-proj-abcdefghijklmnopqrstuvwxyz1234567890",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateAnthropicKey(tt.key)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateAnthropicKey() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateOpenAIKey(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name      string
		key       string
		wantError bool
	}{
		{
			name:      "valid key with sk prefix",
			key:       "sk-abcdefghijklmnopqrstuvwxyz1234567890",
			wantError: false,
		},
		{
			name:      "valid key with sk-proj prefix",
			key:       "sk-proj-abcdefghijklmnopqrstuvwxyz1234567890",
			wantError: false,
		},
		{
			name:      "empty key",
			key:       "",
			wantError: true,
		},
		{
			name:      "too short",
			key:       "sk-abc",
			wantError: true,
		},
		{
			name:      "missing prefix",
			key:       "abcdefghijklmnopqrstuvwxyz1234567890",
			wantError: true,
		},
		{
			name:      "placeholder value",
			key:       "sk-your-api-key",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateOpenAIKey(tt.key)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateOpenAIKey() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateGoogleKey(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name      string
		key       string
		wantError bool
	}{
		{
			name:      "valid key",
			key:       "AIzaSyAbCdEfGhIjKlMnOpQrStUvWxYz1234567",
			wantError: false,
		},
		{
			name:      "empty key",
			key:       "",
			wantError: true,
		},
		{
			name:      "too short",
			key:       "AIzaSy",
			wantError: true,
		},
		{
			name:      "contains spaces",
			key:       "AIzaSyAbCdEf GhIjKlMnOpQrStUvWxYz1234567",
			wantError: true,
		},
		{
			name:      "placeholder value",
			key:       "your-api-key",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateGoogleKey(tt.key)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateGoogleKey() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateAll(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name         string
		anthropicKey string
		openaiKey    string
		googleKey    string
		wantError    bool
	}{
		{
			name:         "all valid",
			anthropicKey: "sk-ant-api03-abcdefghijklmnopqrstuvwxyz1234567890",
			openaiKey:    "sk-abcdefghijklmnopqrstuvwxyz1234567890",
			googleKey:    "AIzaSyAbCdEfGhIjKlMnOpQrStUvWxYz1234567",
			wantError:    false,
		},
		{
			name:         "all empty (valid - no validation needed)",
			anthropicKey: "",
			openaiKey:    "",
			googleKey:    "",
			wantError:    false,
		},
		{
			name:         "one invalid",
			anthropicKey: "invalid-key",
			openaiKey:    "sk-abcdefghijklmnopqrstuvwxyz1234567890",
			googleKey:    "AIzaSyAbCdEfGhIjKlMnOpQrStUvWxYz1234567",
			wantError:    true,
		},
		{
			name:         "multiple invalid",
			anthropicKey: "invalid",
			openaiKey:    "also-invalid",
			googleKey:    "still-invalid",
			wantError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateAll(tt.anthropicKey, tt.openaiKey, tt.googleKey)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateAll() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestMaskKey(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected string
	}{
		{
			name:     "normal key",
			key:      "sk-ant-api03-abcdefghijklmnopqrstuvwxyz1234567890",
			expected: "sk...90",
		},
		{
			name:     "short key",
			key:      "short",
			expected: "sh...rt",
		},
		{
			name:     "empty key",
			key:      "",
			expected: "<empty>",
		},
		{
			name:     "exactly 4 chars",
			key:      "1234",
			expected: "****",
		},
		{
			name:     "5 chars shows masking",
			key:      "12345",
			expected: "12...45",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskKey(tt.key)
			if result != tt.expected {
				t.Errorf("MaskKey() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestValidationError(t *testing.T) {
	err := &ValidationError{
		Provider: "TestProvider",
		Reason:   "test reason",
	}

	expected := "invalid credential for TestProvider: test reason"
	if err.Error() != expected {
		t.Errorf("ValidationError.Error() = %v, want %v", err.Error(), expected)
	}
}

func TestValidateBasic(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name      string
		key       string
		wantError bool
	}{
		{
			name:      "valid key",
			key:       "validkey123456",
			wantError: false,
		},
		{
			name:      "empty key",
			key:       "",
			wantError: true,
		},
		{
			name:      "too short",
			key:       "short",
			wantError: true,
		},
		{
			name:      "too long",
			key:       strings.Repeat("a", 600),
			wantError: true,
		},
		{
			name:      "placeholder - your-api-key",
			key:       "your-api-key",
			wantError: true,
		},
		{
			name:      "placeholder - test",
			key:       "test",
			wantError: true,
		},
		{
			name:      "placeholder - example",
			key:       "example",
			wantError: true,
		},
		{
			name:      "placeholder - xxx",
			key:       "xxx",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateBasic("TestProvider", tt.key)
			if (err != nil) != tt.wantError {
				t.Errorf("validateBasic() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

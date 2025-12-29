package config

import (
	"os"
	"testing"
	"time"
)

func TestGetString(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		envValue     string
		defaultValue string
		expected     string
	}{
		{
			name:         "env var set",
			envKey:       "TEST_STRING",
			envValue:     "hello",
			defaultValue: "default",
			expected:     "hello",
		},
		{
			name:         "env var not set",
			envKey:       "TEST_STRING_UNSET",
			envValue:     "",
			defaultValue: "default",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.envKey, tt.envValue)
				defer os.Unsetenv(tt.envKey)
			}

			got := GetString(tt.envKey, tt.defaultValue)
			if got != tt.expected {
				t.Errorf("GetString() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetInt(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		envValue     string
		defaultValue int
		expected     int
	}{
		{
			name:         "valid int",
			envKey:       "TEST_INT",
			envValue:     "42",
			defaultValue: 10,
			expected:     42,
		},
		{
			name:         "invalid int",
			envKey:       "TEST_INT_INVALID",
			envValue:     "not-a-number",
			defaultValue: 10,
			expected:     10,
		},
		{
			name:         "env var not set",
			envKey:       "TEST_INT_UNSET",
			envValue:     "",
			defaultValue: 10,
			expected:     10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.envKey, tt.envValue)
				defer os.Unsetenv(tt.envKey)
			}

			got := GetInt(tt.envKey, tt.defaultValue)
			if got != tt.expected {
				t.Errorf("GetInt() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetDuration(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		envValue     string
		defaultValue time.Duration
		expected     time.Duration
	}{
		{
			name:         "valid duration",
			envKey:       "TEST_DURATION",
			envValue:     "5s",
			defaultValue: 10 * time.Second,
			expected:     5 * time.Second,
		},
		{
			name:         "invalid duration",
			envKey:       "TEST_DURATION_INVALID",
			envValue:     "not-a-duration",
			defaultValue: 10 * time.Second,
			expected:     10 * time.Second,
		},
		{
			name:         "env var not set",
			envKey:       "TEST_DURATION_UNSET",
			envValue:     "",
			defaultValue: 10 * time.Second,
			expected:     10 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.envKey, tt.envValue)
				defer os.Unsetenv(tt.envKey)
			}

			got := GetDuration(tt.envKey, tt.defaultValue)
			if got != tt.expected {
				t.Errorf("GetDuration() = %v, want %v", got, tt.expected)
			}
		})
	}
}

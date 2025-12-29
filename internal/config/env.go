// Package config provides utilities for environment variable configuration.
package config

import (
	"fmt"
	"os"
	"time"
)

// GetString retrieves a string environment variable or returns the default value if not set.
func GetString(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

// GetInt retrieves an integer environment variable or returns the default value if not set or invalid.
func GetInt(key string, defaultValue int) int {
	if v := os.Getenv(key); v != "" {
		var parsed int
		if n, err := fmt.Sscanf(v, "%d", &parsed); err == nil && n == 1 {
			return parsed
		}
	}
	return defaultValue
}

// GetDuration retrieves a duration environment variable or returns the default value if not set or invalid.
func GetDuration(key string, defaultValue time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if parsed, err := time.ParseDuration(v); err == nil {
			return parsed
		}
	}
	return defaultValue
}

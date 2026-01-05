package config

import (
	"log/slog"
	"os"
)

// GetEnv retrieves the value of the environment variable named by the key.
// It returns the value, which will be the fallback if the variable is not present.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// MustGetEnv retrieves the value of the environment variable named by the key.
// It exits the application if the variable is not present or empty.
func MustGetEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		slog.Error("Missing required environment variable", "key", key)
		os.Exit(1)
	}
	return value
}

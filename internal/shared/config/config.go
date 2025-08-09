package config

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Env           string
	HttpPort      string
	Timezone      string
	// Google
	GCalCalendarID string
	GoogleCredsJSON string // base64 or path
	// Supabase
	SupabaseURL string
	SupabaseKey string
	// Redis
	RedisURL string
}

// Load reads configuration from environment variables
func Load() (Config, error) {
	// Try to load .env file if it exists
	godotenv.Load()
	
	// Also try to load config.env if .env doesn't exist
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		godotenv.Load("config.env")
	}

	config := Config{
		Env:           getEnv("APP_ENV", "development"),
		HttpPort:      getEnv("HTTP_PORT", "8080"),
		Timezone:      getEnv("TZ", "UTC"),
		GCalCalendarID: getEnv("GCAL_CALENDAR_ID", ""),
		GoogleCredsJSON: getEnv("GOOGLE_CREDS_JSON", ""),
		SupabaseURL:   getEnv("SUPABASE_URL", ""),
		SupabaseKey:   getEnv("SUPABASE_KEY", ""),
		RedisURL:      getEnv("REDIS_URL", "redis://localhost:6379"),
	}

	// Validate required fields only in production
	if config.Env == "production" {
		if config.GCalCalendarID == "" {
			return config, fmt.Errorf("GCAL_CALENDAR_ID is required in production")
		}
		if config.GoogleCredsJSON == "" {
			return config, fmt.Errorf("GOOGLE_CREDS_JSON is required in production")
		}
	}

	return config, nil
}

// getEnv retrieves an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetGoogleCreds returns the Google credentials as bytes
func (c *Config) GetGoogleCreds() ([]byte, error) {
	// Check if it's base64 encoded
	if len(c.GoogleCredsJSON) > 0 {
		// Try to decode as base64 first
		if decoded, err := base64.StdEncoding.DecodeString(c.GoogleCredsJSON); err == nil {
			return decoded, nil
		}
		// If not base64, treat as file path
		if _, err := os.Stat(c.GoogleCredsJSON); err == nil {
			return os.ReadFile(c.GoogleCredsJSON)
		}
		// If neither, treat as raw JSON string
		return []byte(c.GoogleCredsJSON), nil
	}
	return nil, fmt.Errorf("no Google credentials provided")
} 
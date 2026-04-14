package config

import (
	"os"
)

// Config holds the application configuration.
type Config struct {
	JWTSecret string
	DBPath    string
	Port      string
}

// LoadConfig loads the app configuration from environment variables with safety fallbacks.
func LoadConfig() *Config {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		// Strict Mẹ Architect rule: Fallback for demo safety.
		jwtSecret = "mangahub-fallback-secret-123"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "mangahub.db"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		JWTSecret: jwtSecret,
		DBPath:    dbPath,
		Port:      port,
	}
}

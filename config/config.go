package config

import (
	"os"
	"strconv"
)

// Config holds the application configuration.
type Config struct {
	JWTSecret      string
	DBPath         string
	Port           string
	TCPPort        string
	MaxTCPClients  int
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

	tcpPort := os.Getenv("TCP_PORT")
	if tcpPort == "" {
		tcpPort = "9090"
	}

	maxClientsStr := os.Getenv("MAX_TCP_CLIENTS")
	maxClients := 100
	if val, err := strconv.Atoi(maxClientsStr); err == nil {
		maxClients = val
	}

	return &Config{
		JWTSecret:     jwtSecret,
		DBPath:        dbPath,
		Port:          port,
		TCPPort:       tcpPort,
		MaxTCPClients: maxClients,
	}
}

package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/user/mangahub/config"
)

func TestLoadConfig(t *testing.T) {
	t.Run("Default Fallbacks", func(t *testing.T) {
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("DB_PATH")
		os.Unsetenv("PORT")

		cfg := config.LoadConfig()
		assert.Equal(t, "mangahub-fallback-secret-123", cfg.JWTSecret)
		assert.Equal(t, "mangahub.db", cfg.DBPath)
		assert.Equal(t, "8080", cfg.Port)
	})

	t.Run("Environment Variables", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "custom-secret")
		os.Setenv("DB_PATH", "test.db")
		os.Setenv("PORT", "9090")
		defer os.Unsetenv("JWT_SECRET")
		defer os.Unsetenv("DB_PATH")
		defer os.Unsetenv("PORT")

		cfg := config.LoadConfig()
		assert.Equal(t, "custom-secret", cfg.JWTSecret)
		assert.Equal(t, "test.db", cfg.DBPath)
		assert.Equal(t, "9090", cfg.Port)
	})
}

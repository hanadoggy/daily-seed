package config_test

import (
	"daily-seed/internal/config"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	// Save current env vars and restore them after tests
	originalMongoURI := os.Getenv("MONGO_URI")
	originalMongoDBName := os.Getenv("MONGO_DB_NAME")
	originalPort := os.Getenv("PORT")
	defer func() {
		os.Setenv("MONGO_URI", originalMongoURI)
		os.Setenv("MONGO_DB_NAME", originalMongoDBName)
		os.Setenv("PORT", originalPort)
	}()

	t.Run("Default values", func(t *testing.T) {
		os.Unsetenv("MONGO_URI")
		os.Unsetenv("MONGO_DB_NAME")
		os.Unsetenv("PORT")

		cfg := config.Load()
		assert.Equal(t, "mongodb://localhost:27017", cfg.MongoURI)
		assert.Equal(t, "daily_seed", cfg.MongoDBName)
		assert.Equal(t, "8080", cfg.Port)
	})

	t.Run("Custom values", func(t *testing.T) {
		os.Setenv("MONGO_URI", "mongodb://test:27017")
		os.Setenv("MONGO_DB_NAME", "test_db")
		os.Setenv("PORT", "9090")

		cfg := config.Load()
		assert.Equal(t, "mongodb://test:27017", cfg.MongoURI)
		assert.Equal(t, "test_db", cfg.MongoDBName)
		assert.Equal(t, "9090", cfg.Port)
	})
}

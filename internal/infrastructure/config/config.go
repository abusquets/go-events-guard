package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
)

// AppConfig struct stores the application configuration
type AppConfig struct {
	ServerPort            int
	ServerUri             string
	MongoDBUri            string
	MongoDB               string
	RedisUrl              string
	RedisPassword         string
	RedisDB               int
	TokenExpiringDuration int
	ApiName               string
	ApiVersion            string
	JWTSecretKey          string
}

var (
	config     *AppConfig // Global variable to hold the singleton instance
	configOnce sync.Once  // Ensures config is initialized only once
)

// getEnv retrieves environment variables or returns a default value if not set
func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// GetAppConfig initializes the configuration (singleton)
func GetAppConfig() (*AppConfig, error) {
	fmt.Println("Loading configuration")

	var err error

	// Ensure the configuration is created only once
	configOnce.Do(func() {
		config = &AppConfig{
			ServerPort: func() int {
				value, _ := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
				return value
			}(),
			ServerUri:     getEnv("SERVER_HOST", "http://localhost:8082"),
			MongoDBUri:    getEnv("MONGODB_URI", ""),
			MongoDB:       getEnv("MONGODB_DB", "eventsguard"),
			RedisUrl:      getEnv("REDIS_URL", "localhost:6379"),
			RedisPassword: getEnv("REDIS_PWD", ""),
			RedisDB: func() int {
				value, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
				return value
			}(),
			TokenExpiringDuration: func() int {
				value, _ := strconv.Atoi(getEnv("TOKEN_EXPIRING_DURATION", "3600"))
				return value
			}(),
			ApiName:      getEnv("API_NAME", "My API"),
			ApiVersion:   getEnv("API_VERSION", "1.0.0"),
			JWTSecretKey: getEnv("JWT_SECRET, ", "-"),
		}

		if config.MongoDBUri == "" {
			err = errors.New("MongoDBUri configuration is missing")
			config = nil
		}

	})

	// Return the singleton instance and any error
	return config, err
}

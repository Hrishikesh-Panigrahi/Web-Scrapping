package config

import (
	"log"
	"os"
	"strconv"
)

// AppConfig holds all configuration for the application
type AppConfig struct {
	ServerPort     string
	Environment    string
	CSVPath        string
	MaxResults     int
	RequestTimeout int // in seconds
	RateLimitRPM   int // requests per minute
	LogLevel       string
	AllowedOrigins []string
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() *AppConfig {
	config := &AppConfig{
		ServerPort:     getEnv("PORT", "8080"),
		Environment:    getEnv("ENVIRONMENT", "development"),
		CSVPath:        getEnv("CSV_PATH", "products.csv"),
		MaxResults:     getEnvAsInt("MAX_RESULTS", 100),
		RequestTimeout: getEnvAsInt("REQUEST_TIMEOUT", 30),
		RateLimitRPM:   getEnvAsInt("RATE_LIMIT_RPM", 5),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		AllowedOrigins: []string{
			getEnv("ALLOWED_ORIGIN_1", "http://localhost:8080"),
			getEnv("ALLOWED_ORIGIN_2", "http://127.0.0.1:8080"),
		},
	}

	// Validate configuration
	if config.ServerPort == "" {
		log.Fatal("PORT environment variable is required")
	}

	log.Printf("Configuration loaded - Environment: %s, Port: %s", config.Environment, config.ServerPort)
	return config
}

// IsProduction returns true if running in production environment
func (c *AppConfig) IsProduction() bool {
	return c.Environment == "production"
}

// IsDevelopment returns true if running in development environment
func (c *AppConfig) IsDevelopment() bool {
	return c.Environment == "development"
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsInt gets environment variable as integer with fallback
func getEnvAsInt(key string, fallback int) int {
	strValue := getEnv(key, "")
	if strValue == "" {
		return fallback
	}

	if value, err := strconv.Atoi(strValue); err == nil {
		return value
	}

	log.Printf("Warning: Invalid integer value for %s, using default %d", key, fallback)
	return fallback
}

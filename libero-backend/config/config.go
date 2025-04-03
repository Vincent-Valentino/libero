package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	JWT         JWTConfig
	Google      OAuthConfig
	Facebook    OAuthConfig
	GitHub      OAuthConfig
	FrontendURL string // Added Frontend URL
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port int
	Host string
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// JWTConfig holds JWT related configuration
type JWTConfig struct {
	Secret    string
	ExpiresIn int
}

// OAuthConfig holds OAuth provider configuration
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// New creates a new Config instance with values from environment variables
func New() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnvAsInt("PORT", 8080),
			Host: getEnv("HOST", "localhost"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "libero"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:    getEnv("JWT_SECRET", "your_secret_key"), // Ensure this is set securely in env
			ExpiresIn: getEnvAsInt("JWT_EXPIRES_IN", 24*60*60), // Default: 24 hours in seconds
		},
		Google: OAuthConfig{
			ClientID:     getEnv("GOOGLE_CLIENT_ID", ""), // Provide actual default or ensure env var is set
			ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""), // Provide actual default or ensure env var is set
			RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),
		},
		Facebook: OAuthConfig{
			ClientID:     getEnv("FACEBOOK_CLIENT_ID", ""), // Provide actual default or ensure env var is set
			ClientSecret: getEnv("FACEBOOK_CLIENT_SECRET", ""), // Provide actual default or ensure env var is set
			RedirectURL:  getEnv("FACEBOOK_REDIRECT_URL", "http://localhost:8080/auth/facebook/callback"),
		},
		GitHub: OAuthConfig{
			ClientID:     getEnv("GITHUB_CLIENT_ID", ""), // Provide actual default or ensure env var is set
			ClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""), // Provide actual default or ensure env var is set
			RedirectURL:  getEnv("GITHUB_REDIRECT_URL", "http://localhost:8080/auth/github/callback"),
		},
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"), // Added Frontend URL loading (default Vite port)
	}
}

// Helper functions to get environment variables
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" { // Check if value is not empty
		return value
	}
	// Consider logging a warning if using default for sensitive values like secrets
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
package config

import (
	"os"
)

// Config holds all application configuration
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Auth     AuthConfig
	AWS      AWSConfig
	Garmin   GarminConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type ServerConfig struct {
	Port string
	Env  string // development, production
}

type AuthConfig struct {
	JWTSecret     string
	TokenDuration int // in hours
}

type AWSConfig struct {
	Region          string
	S3Bucket        string
	AccessKeyID     string
	SecretAccessKey string
}

type GarminConfig struct {
	ConsumerKey    string
	ConsumerSecret string
	CallbackURL    string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "healthuser"),
			Password: getEnv("DB_PASSWORD", "healthpass"),
			DBName:   getEnv("DB_NAME", "health_assistant"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Auth: AuthConfig{
			JWTSecret:     getEnv("JWT_SECRET", "change-me-in-production"),
			TokenDuration: 24, // 24 hours
		},
		AWS: AWSConfig{
			Region:          getEnv("AWS_REGION", "us-east-1"),
			S3Bucket:        getEnv("AWS_S3_BUCKET", "health-assistant-photos"),
			AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
		},
		Garmin: GarminConfig{
			ConsumerKey:    getEnv("GARMIN_CONSUMER_KEY", ""),
			ConsumerSecret: getEnv("GARMIN_CONSUMER_SECRET", ""),
			CallbackURL:    getEnv("GARMIN_CALLBACK_URL", "http://localhost:8080/auth/garmin/callback"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

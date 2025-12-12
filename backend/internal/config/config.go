package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Environment string
	Port        string

	// Database
	DatabaseURL string

	// Redis
	RedisURL string

	// JWT
	JWTSecret          string
	JWTExpiryHours     int
	RefreshExpiryHours int

	// CORS
	CORSOrigins string

	// Email (for future use)
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	MailFrom     string
}

func Load() *Config {
	// Load .env file if exists
	godotenv.Load()

	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "3000"),

		DatabaseURL: getEnv("DATABASE_URL", "postgres://linkmy:linkmy_password@localhost:5432/linkmy_db?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),

		JWTSecret:          getEnv("JWT_SECRET", "your-super-secret-key-change-in-production"),
		JWTExpiryHours:     24,
		RefreshExpiryHours: 168, // 7 days

		CORSOrigins: getEnv("CORS_ORIGINS", "http://localhost:3001,http://localhost:5173"),

		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		MailFrom:     getEnv("MAIL_FROM", "noreply@linkmy.deepkernel.site"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) IsDevelopment() bool {
	return strings.ToLower(c.Environment) == "development"
}

func (c *Config) IsProduction() bool {
	return strings.ToLower(c.Environment) == "production"
}

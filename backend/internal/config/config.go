package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port       string
	LogLevel   string
	LogFormat  string
	CorsOrigin string
	
}

type S3Config struct {
	Bucket   string
	Endpoint string
	Region   string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func LoadConfig() *Config {
	return &Config{
		Port:      GetEnv("PORT", "8080"),
		LogLevel:  GetEnv("LOG_LEVEL", "info"),
		LogFormat: GetEnv("LOG_FORMAT", "json"),
		CorsOrigin: GetEnv("CORS_ORIGIN", "http://localhost:5173"),
		
		DB: DBConfig{
			Host:     GetEnv("POSTGRES_HOST", "localhost"),
			Port:     GetEnv("POSTGRES_PORT", "5432"),
			User:     GetEnv("POSTGRES_USER", "postgres"),
			Password: GetEnv("POSTGRES_PASSWORD", "postgres"),
			Name:     GetEnv("POSTGRES_DB", "pipeline_processing"),
			SSLMode:  GetEnv("POSTGRES_SSLMODE", "disable"),
		},
	}
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

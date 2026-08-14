package config

import "os"

type Config struct {
	Port string
}

func LoadConfig() *Config {
	port := GetEnv("PORT", "8080")
	return &Config{
		Port: port,
	}
}

func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
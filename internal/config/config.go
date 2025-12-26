package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port  string
	DBUrl string
}

func LoadConfig() *Config {
	// Load .env file if it exists (great for local dev)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	return &Config{
		Port:  getEnv("PORT", ":8080"),
		DBUrl: getEnv("DATABASE_URL", "postgres://postgres:root@localhost:5432/app_db?sslmode=disable"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

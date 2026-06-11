package pkg

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Environment string
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env file: %w", err)
	}

	cfg := Config{
		DatabaseURL: getRequired("DATABASE_URL"),
		Environment: getRequired("ENVIRONMENT"),
	}

	return cfg
}

func getRequired(key string) string {
	value := os.Getenv(key)

	if value == "" {
		log.Fatalf("environment variable %s is required", key)
	}

	return value
}

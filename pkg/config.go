package pkg

import (
	"log"
	"os"
)

type Config struct {
	DatabaseURL string
}

func Load() Config {
	cfg := Config{
		DatabaseURL: getRequired("DATABASE_URL"),
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

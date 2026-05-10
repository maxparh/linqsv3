package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL         string
	JWTSecret     string
	JWTExpiration int // в секундах
	ServerPort    string
	FrontendURL   string // для CORS и генерации полных ссылок
}

func Load() *Config {
	// Загружаем .env файл, если он есть.
	// В Docker-контейнере этого файла может не быть, и это нормально - переменные придут из среды.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it, using system environment variables")
	}

	return &Config{
		DBURL:         getEnvRequired("DATABASE_URL"),
		JWTSecret:     getEnvRequired("JWT_SECRET"),
		JWTExpiration: getEnvAsInt("JWT_EXPIRATION", 3600),
		ServerPort:    getEnv("SERVER_PORT", ":8080"),
		FrontendURL:   getEnv("FRONTEND_URL", "http://localhost:5173"),
	}
}

func getEnvRequired(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		log.Fatalf("Critical environment variable %s is not set", key)
	}
	return value
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if valueStr, ok := os.LookupEnv(key); ok {
		var value int
		if _, err := fmt.Sscanf(valueStr, "%d", &value); err == nil {
			return value
		}
		log.Printf("Invalid integer value for %s, using fallback %d", key, fallback)
	}
	return fallback
}

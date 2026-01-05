package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Host        string
	Port        string
	Username    string
	Password    string
	DBname      string
	SSLMode     string
	MaxAttempts int
}

func LoadConfig() Config {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Проверяем обязательные переменные
	requiredVars := []string{"DB_HOST", "DB_PORT", "DB_USERNAME", "DB_PASSWORD", "DB_NAME"}
	for _, key := range requiredVars {
		if os.Getenv(key) == "" {
			log.Fatalf("Required environment variable %s is not set", key)
		}
	}

	// Читаем MaxAttempts с значением по умолчанию (не критичный параметр)
	maxAttempts := 5
	if val := os.Getenv("DB_MAX_ATTEMPTS"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			maxAttempts = parsed
		}
	}

	// SSLMode с разумным значением по умолчанию
	sslMode := os.Getenv("DB_SSLMODE")
	if sslMode == "" {
		sslMode = "disable"
	}

	return Config{
		Host:        os.Getenv("DB_HOST"),
		Port:        os.Getenv("DB_PORT"),
		Username:    os.Getenv("DB_USERNAME"),
		Password:    os.Getenv("DB_PASSWORD"),
		DBname:      os.Getenv("DB_NAME"),
		SSLMode:     sslMode,
		MaxAttempts: maxAttempts,
	}
}

// mustGetEnv получает значение переменной окружения или паникует
func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Required environment variable %s is not set", key))
	}
	return value
}
